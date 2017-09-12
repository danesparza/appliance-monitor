// +build !linux !arm

package network

import (
	"log"
)

// ResetHostname changes the hostname for the local machine
func UpdateWifiCredentials(ssid, password string) error {
	log.Println("[INFO] Not running on Linux/ARM, so wifi can't be changed...")

	return nil
}
