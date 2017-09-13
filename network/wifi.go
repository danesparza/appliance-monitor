// +build !linux !arm

package network

import "log"

// UpdateWifiCredentials changes the hostname for the local machine
func UpdateWifiCredentials(ssid, password string, reboot chan bool) error {
	log.Printf("[INFO] Not running on Linux/ARM.  Updating the wifi credentials.\nNew SSID: %v\nNew Password: %v\n", ssid, password)

	log.Println("[INFO] Requesting a reboot because of wifi changes")
	reboot <- true

	return nil
}
