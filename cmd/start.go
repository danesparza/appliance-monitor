package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	/*
		Raspberry pi specific imports:
		"github.com/danesparza/appliance-monitor/api"
		"github.com/kidoman/embd"
		_ "github.com/kidoman/embd/host/rpi" // This loads the RPi driver
	*/

	"github.com/danesparza/appliance-monitor/api"
	"github.com/gorilla/mux"
	"github.com/kidoman/embd"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	serverInterface   string
	serverPort        int
	serverUIDirectory string
	allowedOrigins    string

	maxPoints             = 120
	applianceRunThreshold = float64(8)
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the appliance monitor server",
	Long: `Appliance monitor provides its own webserver which can serve both the 
API and the UI for the app in addition to monitoring the system sensors
and providing notification.`,
	Run: serve,
}

func serve(cmd *cobra.Command, args []string) {

	//	If we have a config file, report it:
	if viper.ConfigFileUsed() != "" {
		log.Println("[INFO] Using config file:", viper.ConfigFileUsed())
	}

	log.Printf("[INFO] Config database: %s\n", viper.GetString("datastore.config"))
	log.Printf("[INFO] Activities database: %s\n", viper.GetString("datastore.activity"))

	//	Create a router and setup our REST endpoints...
	var Router = mux.NewRouter()

	//	Setup our routes
	// Router.HandleFunc("/", api.ShowUI)

	//	Activities
	Router.HandleFunc("/activities/get", api.GetActivity).Methods("GET", "POST")

	//	Config
	Router.HandleFunc("/config/getall", nil)
	Router.HandleFunc("/config/get", nil)
	Router.HandleFunc("/config/set", nil)

	//	System information
	Router.HandleFunc("/system/state", api.GetCurrentState).Methods("GET")
	Router.HandleFunc("/system/wifi", nil)

	//	System resets
	Router.HandleFunc("/reset/network", nil)
	Router.HandleFunc("/reset/config", nil)

	//	Websocket connections
	// Router.Handle("/ws", api.WsHandler{H: api.WsHub})

	//	Trap program exit appropriately
	ctx, cancel := context.WithCancel(context.Background())
	sigch := make(chan os.Signal, 2)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	go handleSignals(ctx, sigch, cancel)

	/*
		//	Start the collection ticker
		go func() {
			log.Println("Initializing GPIO...")
			embd.InitGPIO()
			defer embd.CloseGPIO()

			pin, err := embd.NewDigitalPin("GPIO_4")
			if err != nil {
				log.Fatal("opening pin:", err)
			}
			defer resetPin(pin)

			if err = pin.SetDirection(embd.Out); err != nil {
				log.Fatal("setting pin direction:", err)
			}

			log.Println("Initializing I2C...")
			if err := embd.InitI2C(); err != nil {
				panic(err)
			}
			defer embd.CloseI2C()

			bus := embd.NewI2CBus(1)
			sensor := envirophat.New(bus)

			log.Println("Initializing InfluxDB client...")
			c, err := client.NewHTTPClient(client.HTTPConfig{Addr: "http://chile.lan:8086"})
			if err != nil {
				log.Fatal(err)
			}

			// Pushover client settings
			log.Println("Initializing Pushover client...")
			pushClient := pushover.New("ad2ujxv7zi7i5zw8fuvt5hu3chjuv4")
			recipient := pushover.NewRecipient("gqukgkJyLtchaLE41WUEJ2qFM7Q3tb")

			hostname, _ := os.Hostname()
			log.Printf("Using hostname %v...", hostname)

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

					// Create a new point batch
					bp, err := client.NewBatchPoints(client.BatchPointsConfig{
						Database:  "sensors",
						Precision: "s",
					})
					if err != nil {
						log.Fatal(err)
					}

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

					// Create a point and add to batch
					tags := map[string]string{"host": hostname}
					fields := map[string]interface{}{
						"x":    x,
						"y":    y,
						"z":    z,
						"xdev": xdevcurrent,
						"ydev": ydevcurrent,
						"zdev": zdevcurrent}

					pt, err := client.NewPoint("envirophat-lsm303d", tags, fields, time.Now())
					if err != nil {
						log.Fatal(err)
					}
					bp.AddPoint(pt)

					// Write the batch
					if err := c.Write(bp); err != nil {
						log.Fatal(err)
					}

					//	Calculate ... are we currently running?
					if ((xdevcurrent * 1000) > applianceRunThreshold) && ((ydevcurrent * 1000) > applianceRunThreshold) && ((zdevcurrent * 1000) > applianceRunThreshold) && currentlyRunning == false {
						log.Println("Looks like the machine is running")
						currentlyRunning = true
						timeStart = time.Now()
					}

					if ((xdevcurrent * 1000) < applianceRunThreshold) && ((ydevcurrent * 1000) < applianceRunThreshold) && ((zdevcurrent * 1000) < applianceRunThreshold) && currentlyRunning == true {
						log.Println("Looks like the machine is stopped")
						currentlyRunning = false

						//	Calculate the run time:
						runningTime := time.Since(timeStart)

						// Send the message to the recipient
						message := pushover.NewMessage(fmt.Sprintf("The dryer has finished running.  It ran for about %v minutes", int(runningTime.Minutes())))
						message.Sound = "bike"
						_, err := pushClient.SendMessage(message, recipient)
						if err != nil {
							log.Printf("Error sending pushover message: %v\n", err)
						}
					}

					pin.Write(embd.Low)

				}
			}
		}()
	*/

	//	If we don't have a UI directory specified...
	/*
		if viper.GetString("server.ui-dir") == "" {
			//	Use the static assets file generated with
			//	https://github.com/elazarl/go-bindata-assetfs using the application-monitor-ui from
			//	https://github.com/danesparza/application-monitor-ui.
			//
			//	To generate this file, place the 'ui'
			//	directory under the main application-monitor directory and run the commands:
			//	go-bindata-assetfs.exe -pkg cmd ./ui/...
			//	mv bindata_assetfs.go cmd
			//	go install ./...
			Router.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(assetFS())))
		} else {
			//	Use the supplied directory:
			log.Printf("[INFO] Using UI directory: %s\n", viper.GetString("server.ui-dir"))
			Router.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir(viper.GetString("server.ui-dir")))))
		}
	*/

	//	Setup the CORS options:
	log.Printf("[INFO] Allowed CORS origins: %s\n", viper.GetString("server.allowed-origins"))
	c := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(allowedOrigins, ","),
		AllowCredentials: true,
	})
	corsHandler := c.Handler(Router)

	//	Format the bound interface:
	formattedInterface := viper.GetString("server.bind")
	if formattedInterface == "" {
		formattedInterface = "127.0.0.1"
	}

	//	If we have an SSL cert specified, use it:
	if viper.GetString("server.sslcert") != "" {
		log.Printf("[INFO] Using SSL cert: %s\n", viper.GetString("server.sslcert"))
		log.Printf("[INFO] Using SSL key: %s\n", viper.GetString("server.sslkey"))
		log.Printf("[INFO] Starting HTTPS server: https://%s:%s\n", formattedInterface, viper.GetString("server.port"))

		log.Printf("[ERROR] %v\n", http.ListenAndServeTLS(viper.GetString("server.bind")+":"+viper.GetString("server.port"), viper.GetString("server.sslcert"), viper.GetString("server.sslkey"), corsHandler))
	} else {
		log.Printf("[INFO] Starting HTTP server: http://%s:%s\n", formattedInterface, viper.GetString("server.port"))
		log.Printf("[ERROR] %v\n", http.ListenAndServe(viper.GetString("server.bind")+":"+viper.GetString("server.port"), corsHandler))
	}
}

func init() {
	RootCmd.AddCommand(startCmd)

	//	Setup our flags
	startCmd.Flags().IntVarP(&serverPort, "port", "p", 1313, "port on which the server will listen")
	startCmd.Flags().StringVarP(&serverInterface, "bind", "i", "", "interface to which the server will bind")
	startCmd.Flags().StringVarP(&serverUIDirectory, "ui-dir", "u", "", "directory for the UI")
	startCmd.Flags().StringVarP(&allowedOrigins, "allowed-origins", "o", "", "comma seperated list of allowed CORS origins")

	//	Bind config flags for optional config file override:
	viper.BindPFlag("server.port", startCmd.Flags().Lookup("port"))
	viper.BindPFlag("server.bind", startCmd.Flags().Lookup("bind"))
	viper.BindPFlag("server.ui-dir", startCmd.Flags().Lookup("ui-dir"))
	viper.BindPFlag("server.allowed-origins", startCmd.Flags().Lookup("allowed-origins"))

	//	Set the build version information:
	api.BuildVersion = BuildVersion
	api.CommitID = CommitID

}

func handleSignals(ctx context.Context, sigch <-chan os.Signal, cancel context.CancelFunc) {
	select {
	case <-ctx.Done():
	case sig := <-sigch:
		switch sig {
		case os.Interrupt:
			log.Println("[INFO] SIGINT")
		case syscall.SIGTERM:
			log.Println("[INFO] SIGTERM")
		}
		log.Println("[INFO] Shutting down ...")
		cancel()
		os.Exit(0)
	}
}

func resetPin(pin embd.DigitalPin) {
	if err := pin.SetDirection(embd.In); err != nil {
		log.Fatal("resetting pin:", err)
	}
	pin.Close()
}
