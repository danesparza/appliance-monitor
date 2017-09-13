// +build !linux !arm

package network

import (
	"log"
	"os"
)

// ResetHostname changes the hostname for the local machine
func ResetHostname(newname string, reboot chan bool) error {
	log.Println("[INFO] Not running on Linux/ARM, so hostname won't get reset...")

	//	Get hostname:
	name, err := os.Hostname()
	if err != nil {
		log.Printf("[ERROR] Problem getting hostname: %v", err.Error())
	}

	log.Printf("[INFO] Current hostname: %v -- desired new hostname: %v\n", name, newname)

	//	Indicate we should trigger a reboot
	log.Println("[INFO] Requesting a reboot because of hostname changes")
	reboot <- true

	return nil
}
