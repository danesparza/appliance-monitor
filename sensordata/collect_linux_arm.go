package sensordata

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/danesparza/appliance-monitor/data"
	"github.com/gregdel/pushover"
	influxdb "github.com/influxdata/influxdb/client/v2"
	"github.com/montanaflynn/stats"
	"github.com/spf13/viper"

	"github.com/danesparza/embd/sensor/envirophat"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
)

// CollectAndProcess performs the sensor data collection and data processing
func CollectAndProcess(ctx context.Context) {
	//	Connect to the datastores:
	log.Printf("[INFO] Config database: %s\n", viper.GetString("datastore.config"))
	configDB := data.ConfigDB{
		Database: viper.GetString("datastore.config")}
	log.Printf("[INFO] Activities database: %s\n", viper.GetString("datastore.activity"))
	activityDB := data.ActivityDB{
		Database: viper.GetString("datastore.activity")}

	//	Initialize GPIO / I2C / etc
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
			//	Turn the LED on
			pin.Write(embd.High)

			/**************************
				SENSOR DATA GATHERING
			***************************/
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

			/***************************
				INFLUXDB DEBUGGING
			****************************/
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

			/*********************************
				EVENT: MACHINE STARTED
			**********************************/
			if ((xdevcurrent * 1000) > applianceRunThreshold) && ((ydevcurrent * 1000) > applianceRunThreshold) && ((zdevcurrent * 1000) > applianceRunThreshold) && currentlyRunning == false {
				log.Println("[DEBUG] Looks like the machine is running")
				WsHub.Broadcast <- []byte("Appliance state: running")

				currentlyRunning = true
				timeStart = time.Now()

				//	Track the activity:
				newActivity := data.Activity{Type: data.ApplianceRunning, Timestamp: time.Now()}
				activityDB.Add(newActivity)
				trackActivity(newActivity, configDB)
			}

			/*********************************
				EVENT: MACHINE STOPPED
			**********************************/
			if ((xdevcurrent * 1000) < applianceRunThreshold) && ((ydevcurrent * 1000) < applianceRunThreshold) && ((zdevcurrent * 1000) < applianceRunThreshold) && currentlyRunning == true {
				log.Println("[DEBUG] Looks like the machine is stopped")
				WsHub.Broadcast <- []byte("Appliance state: stopped")

				currentlyRunning = false

				//	Calculate the run time:
				runningTime := time.Since(timeStart)

				//	Track the activity:
				newActivity := data.Activity{Type: data.ApplianceStopped, Timestamp: time.Now()}
				activityDB.Add(newActivity)
				trackActivity(newActivity, configDB)

				// Send a Pushover message
				err := sendPushoverNotification(configDB, int(runningTime.Minutes()))
				if err != nil {
					log.Printf("[WARN] Problem sending pushover message: %v\n", err)
				}

			}

			//	Turn the LED off
			pin.Write(embd.Low)
		}
	}
}

// Resets the pin
func resetPin(pin embd.DigitalPin) {
	if err := pin.SetDirection(embd.In); err != nil {
		log.Fatal("resetting pin:", err)
	}
	pin.Close()
}

// Track the activity in the cloud
func trackActivity(activity data.Activity, c data.ConfigDB) error {
	url := "https://api.appliance-monitor.com/v1/activity"

	//	Get the deviceID:
	deviceID, _ := c.Get("deviceID")

	//	Serialize to JSON
	cloudActivity := data.CloudActivity{
		DeviceID:  deviceID.Value,
		Timestamp: activity.Timestamp,
		Type:      activity.Type}

	jsonStr, err := json.Marshal(cloudActivity)
	if err != nil {
		log.Printf("[WARN] Problem marshalling cloud activity message: %v\n", err)
	}

	//	Create a request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	//	Set our headers
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	//	Execute the request
	client := &http.Client{}
	_, err = client.Do(req)

	return err
}

// Send a pushover notification
func sendPushoverNotification(c data.ConfigDB, runningTime int) error {

	//	Get the config data
	pushAPIkey, err := c.Get("pushoverapikey")
	pushTo, err := c.Get("pushoverrecipient")
	applianceName, err := c.Get("name")

	//	If we have config data set...
	if err == nil && pushTo.Value != "" {
		//	Create a new client and push a message
		pushClient := pushover.New(pushAPIkey.Value)
		recipient := pushover.NewRecipient(pushTo.Value)
		message := pushover.NewMessage(fmt.Sprintf("%v has finished running.  It ran for about %v minutes", applianceName, runningTime))
		message.Sound = "bike"
		_, err := pushClient.SendMessage(message, recipient)
		if err != nil {
			return err
		}
	}

	return err
}
