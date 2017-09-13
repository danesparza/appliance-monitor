package network

import (
	"io/ioutil"
	"log"
	"syscall"
	"time"
)

// UpdateWifiCredentials updates the ssid and password used
func UpdateWifiCredentials(ssid, password string) error {
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

	RebootMachine()

	return nil
}

// RebootMachine calls sync and reboots the machine
func RebootMachine() {
	log.Println("[INFO] Rebooting...")
	syscall.Sync()

	// Wait before exiting, in order to give our parent enough time to finish
	countdownBeforeExit := time.NewTimer(time.Second * 3)
	<-countdownBeforeExit.C

	syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
}
