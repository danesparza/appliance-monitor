// +build !linux !arm

package sensordata

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

// CollectAndProcess performs the sensor data collection and data processing
func CollectAndProcess(ctx context.Context) {
	log.Println("[INFO] Running on a platform other than Linux/ARM, so this will be boring...")

	//	Connect to the datastores:
	log.Printf("[INFO] Config database: %s\n", viper.GetString("datastore.config"))
	log.Printf("[INFO] Activities database: %s\n", viper.GetString("datastore.activity"))

	hostname, _ := os.Hostname()
	log.Printf("[INFO] Using hostname %v...", hostname)

	//	Keep track of state of device
	currentlyRunning := false

	//	Loop and respond to channels:
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(1 * time.Second):

			//	Dummy activity:
			if !currentlyRunning {
				//	We should actually log to the activity datastore:
				WsHub.Broadcast <- []byte("Appliance state: running")
				currentlyRunning = true
			} else {
				WsHub.Broadcast <- []byte("Appliance state: stopped")
				currentlyRunning = false
			}

		}
	}
}
