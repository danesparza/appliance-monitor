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
	log.Println("[INFO] Rebooting...")

	// Wait before exiting, in order to give our parent enough time to finish
	countdownBeforeExit := time.NewTimer(time.Second * 3)
	<-countdownBeforeExit.C

	log.Println("[WARN] Not running on Linux/ARM, so not rebooting")
}
