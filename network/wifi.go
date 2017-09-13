// +build !linux !arm

package network

import (
	"log"
	"time"
)

// UpdateWifiCredentials changes the hostname for the local machine
func UpdateWifiCredentials(ssid, password string) error {
	return nil // fmt.Errorf("Not running on Linux/ARM, so wifi can't be changed")
}

// RebootMachine calls sync and reboots the machine
func RebootMachine() {
	go func() {
		log.Println("[INFO] Rebooting...")
	}()

	// Wait before exiting, in order to give our parent enough time to finish
	time.Sleep(3 * time.Second)

	log.Println("[WARN] Not running on Linux/ARM, so not rebooting")
}
