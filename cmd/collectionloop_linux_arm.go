package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/danesparza/appliance-monitor/data"
	"github.com/gregdel/pushover"
	influxdb "github.com/influxdb/influxdb/client/v2"
	"github.com/montanaflynn/stats"
	"github.com/spf13/viper"

	"github.com/danesparza/embd/sensor/envirophat"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
)

func collectionloop(ctx context.Context) {
	//	Connect to the datastores:
	log.Printf("[INFO] Config database: %s\n", viper.GetString("datastore.config"))
	configDB := data.ConfigDB{
		Database: viper.GetString("datastore.config")}
	log.Printf("[INFO] Activities database: %s\n", viper.GetString("datastore.activity"))
	activityDB := data.ActivityDB{
		Database: viper.GetString("datastore.activity")}

	log.Println("[INFO] Initializing GPIO...")
	err := embd.InitGPIO()
	if err != nil {
		log.Fatal("[ERROR] Initializing GPIO:", err)
	}
	defer embd.CloseGPIO()

	pin, err := embd.NewDigitalPin("GPIO_4")
	if err != nil {
		log.Fatal("[ERROR] opening pin:", err)
	}
	defer resetPin(pin)

	if err = pin.SetDirection(embd.Out); err != nil {
		log.Fatal("[ERROR] setting pin direction:", err)
	}

	log.Println("[INFO] Initializing I2C...")
	if err := embd.InitI2C(); err != nil {
		log.Fatal("[ERROR] Initalizing I2C:", err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)
	sensor := envirophat.New(bus)

	hostname, _ := os.Hostname()
	log.Printf("[INFO] Using hostname %v...", hostname)

	//	influxserver should be a url, like
	//	http://chile.lan:8086
	influxURL, err := configDB.Get("influxserver")
	if err == nil && influxURL.Value != "" {
		log.Printf("[INFO] Using influxserver %v...", influxURL.Value)
	}
	c, _ := influxdb.NewHTTPClient(influxdb.HTTPConfig{Addr: influxURL.Value})

	//	Keep track of the axis values
	var xaxis, yaxis, zaxis []float64
	var xdev, ydev, zdev []float64

	//	Keep track of state of device
	currentlyRunning := false
	timeStart := time.Now()

	//	Loop and respond to channels:
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(1 * time.Second):
			//	Perform sensor data gathering
			pin.Write(embd.High)

			//	Get accelerometer values from the sensor
			x, y, z, err := sensor.Accelerometer()
			if err != nil {
				log.Fatal(err)
			}

			//	Track the measurements
			xaxis = append(xaxis, x)
			yaxis = append(yaxis, y)
			zaxis = append(zaxis, z)

			//	Calculate standard deviation
			xdevcurrent, err := stats.StandardDeviation(xaxis)
			if err != nil {
				log.Fatal(err)
			}
			xdev = append(xdev, xdevcurrent)

			ydevcurrent, err := stats.StandardDeviation(yaxis)
			if err != nil {
				log.Fatal(err)
			}
			ydev = append(ydev, ydevcurrent)

			zdevcurrent, err := stats.StandardDeviation(zaxis)
			if err != nil {
				log.Fatal(err)
			}
			zdev = append(zdev, zdevcurrent)

			//	Keep a rolling collection of data...
			//	If we already have maxPoints items
			//	remove the first item:
			if len(xaxis) > maxPoints {
				xaxis = xaxis[1:]
			}

			if len(yaxis) > maxPoints {
				yaxis = yaxis[1:]
			}

			if len(zaxis) > maxPoints {
				zaxis = zaxis[1:]
			}

			if len(zdev) > maxPoints {
				xdev = xdev[1:]
			}

			if len(ydev) > maxPoints {
				ydev = ydev[1:]
			}

			if len(zdev) > maxPoints {
				zdev = zdev[1:]
			}

			if influxURL.Value != "" {
				// Create a new point batch
				bp, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
					Database:  "sensors",
					Precision: "s",
				})
				if err != nil {
					log.Fatal(err)
				}

				// Create a point and add to batch
				tags := map[string]string{"host": hostname}
				fields := map[string]interface{}{
					"x":    x,
					"y":    y,
					"z":    z,
					"xdev": xdevcurrent,
					"ydev": ydevcurrent,
					"zdev": zdevcurrent}

				pt, err := influxdb.NewPoint("envirophat-lsm303d", tags, fields, time.Now())
				if err != nil {
					log.Fatal(err)
				}
				bp.AddPoint(pt)

				// Write the batch
				if err := c.Write(bp); err != nil {
					log.Printf("[WARN] Problem writing to InfluxDB server: %v", err)
				}
			}

			//	Calculate ... are we currently running?
			if ((xdevcurrent * 1000) > applianceRunThreshold) && ((ydevcurrent * 1000) > applianceRunThreshold) && ((zdevcurrent * 1000) > applianceRunThreshold) && currentlyRunning == false {
				log.Println("[DEBUG] Looks like the machine is running")

				currentlyRunning = true
				timeStart = time.Now()

				//	Track the activity:
				activityDB.Add(data.Activity{Type: data.ApplianceRunning})
			}

			if ((xdevcurrent * 1000) < applianceRunThreshold) && ((ydevcurrent * 1000) < applianceRunThreshold) && ((zdevcurrent * 1000) < applianceRunThreshold) && currentlyRunning == true {
				log.Println("[DEBUG] Looks like the machine is stopped")
				currentlyRunning = false

				//	Calculate the run time:
				runningTime := time.Since(timeStart)

				//	Track the activity:
				activityDB.Add(data.Activity{Type: data.ApplianceStopped})

				// Send a Pushover message
				pushAPIkey, err := configDB.Get("pushoverapikey")
				pushTo, err := configDB.Get("pushoverrecipient")
				if err == nil && pushTo.Value != "" {
					pushClient := pushover.New(pushAPIkey.Value)
					recipient := pushover.NewRecipient(pushTo.Value)
					message := pushover.NewMessage(fmt.Sprintf("The dryer has finished running.  It ran for about %v minutes", int(runningTime.Minutes())))
					message.Sound = "bike"
					_, err := pushClient.SendMessage(message, recipient)
					if err != nil {
						log.Printf("[WARN] Problem sending pushover message: %v\n", err)
					}
				}
			}

			pin.Write(embd.Low)

		}
	}
}

func resetPin(pin embd.DigitalPin) {
	if err := pin.SetDirection(embd.In); err != nil {
		log.Fatal("resetting pin:", err)
	}
	pin.Close()
}
