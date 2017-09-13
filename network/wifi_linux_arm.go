package network

import (
	"io/ioutil"
	"log"
)

// UpdateWifiCredentials updates the ssid and password used
func UpdateWifiCredentials(ssid, password string, reboot chan bool) error {
	log.Printf("[INFO] Updating the wifi credentials.\nNew SSID: %v\nNew Password: %v\n", ssid, password)

	//	Get the formatted config file
	formattedConfig, err := FormatWifiCredentials(ssid, password)
	if err != nil {
		log.Printf("[ERROR] Formatting wifi credentials: %v", err.Error())
		return err
	}

	//	Save the formatted config file
	err = ioutil.WriteFile("/etc/wpa_supplicant/wpa_supplicant.conf", []byte(formattedConfig), 0600)
	if err != nil {
		log.Printf("[ERROR] Problem writing /etc/wpa_supplicant/wpa_supplicant.conf: %v", err.Error())
		return err
	}

	log.Println("[INFO] Requesting a reboot because of wifi changes")
	reboot <- true

	return nil
}
