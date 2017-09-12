package network

import (
	"log"
)

// UpdateWifiCredentials updates the ssid and password used
func UpdateWifiCredentials(ssid, password string) error {
	log.Printf("[INFO] Updating the wifi credentials.\nNew SSID: %v\nNew Password: %v\n", ssid, password)

	return nil
}
