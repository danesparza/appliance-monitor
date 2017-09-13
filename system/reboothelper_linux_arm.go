package system

import (
	"log"
	"syscall"
	"time"
)

// ListenForReboots listens for requests on the 'reboot' channel and reboots the system
func ListenForReboots(reboot chan bool) {

	//	Loop and listen for requests on the 'reboot' channel
	for {
		select {
		case <-reboot:
			//	Sync filesystem
			syscall.Sync()

			//	Settle a bit before rebooting...
			time.Sleep(3 * time.Second)
			log.Println("[INFO] Rebooting...")

			// Wait before exiting, in order to give our parent enough time to finish
			time.Sleep(2 * time.Second)
			syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
		}
	}
}
