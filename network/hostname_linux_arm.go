package network

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// ResetHostname changes the hostname for the local machine
func ResetHostname(newname string) error {
	//	Get hostname:
	name, err := os.Hostname()
	if err != nil {
		log.Printf("[ERROR] Problem getting hostname: %v", err.Error())
		return err
	}

	log.Printf("[INFO] Current hostname: %v -- desired new hostname: %v\n", name, newname)

	//	Update /etc/hostname file
	hostname, err := ioutil.ReadFile("/etc/hostname")
	if err != nil {
		log.Printf("[ERROR] Problem reading /etc/hostname: %v", err.Error())
		return err
	}
	newhostname := strings.Replace(string(hostname[:]), "raspberrypi", newname, -1)
	err = ioutil.WriteFile("/etc/hostname", newhostname, 0644)
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
	newhosts := strings.Replace(string(hosts[:]), "raspberrypi", newname, -1)
	err = ioutil.WriteFile("/etc/hosts", newhosts, 0644)
	if err != nil {
		log.Printf("[ERROR] Problem writing /etc/hosts: %v", err.Error())
		return err
	}

	return nil
}
