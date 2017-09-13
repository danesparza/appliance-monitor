package system

import (
	"context"
	"log"
	"syscall"
	"time"
)

// ListenForReboots listens for requests on the 'reboot' channel and reboots the system
func ListenForReboots(ctx context.Context, reboot chan bool) {

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
		case <-ctx.Done():
			log.Println("[INFO] Stopping the reboot helper")
			return
		}
	}
}
