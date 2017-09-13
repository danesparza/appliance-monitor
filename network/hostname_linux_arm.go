package network

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// ResetHostname changes the hostname for the local machine
func ResetHostname(newname string, reboot chan bool) error {
	//	Get hostname:
	currentname, err := os.Hostname()
	if err != nil {
		log.Printf("[ERROR] Problem getting hostname: %v", err.Error())
		return err
	}

	log.Printf("[INFO] Current hostname: %v -- desired new hostname: %v\n", currentname, newname)

	//	Update /etc/hostname file
	hostname, err := ioutil.ReadFile("/etc/hostname")
	if err != nil {
		log.Printf("[ERROR] Problem reading /etc/hostname: %v", err.Error())
		return err
	}
	newhostname := strings.Replace(string(hostname[:]), currentname, newname, -1)
	err = ioutil.WriteFile("/etc/hostname", []byte(newhostname), 0644)
	if err != nil {
		log.Printf("[ERROR] Problem writing /etc/hostname: %v", err.Error())
		return err
	}

	//	Update the /etc/hosts file
	hosts, err := ioutil.ReadFile("/etc/hosts")
	if err != nil {
		log.Printf("[ERROR] Problem reading /etc/hosts: %v", err.Error())
		return err
	}
	newhosts := strings.Replace(string(hosts[:]), currentname, newname, -1)
	err = ioutil.WriteFile("/etc/hosts", []byte(newhosts), 0644)
	if err != nil {
		log.Printf("[ERROR] Problem writing /etc/hosts: %v", err.Error())
		return err
	}

	//	Indicate we should trigger a reboot
	log.Println("[INFO] Requesting a reboot because of hostname changes")
	reboot <- true

	return nil
}
