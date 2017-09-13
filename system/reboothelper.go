// +build !linux !arm

package system

import (
	"context"
	"log"
	"time"
)

// ListenForReboots listens for requests on the 'reboot' channel and reboots the system
func ListenForReboots(ctx context.Context, reboot chan bool) {

	//	Loop and listen for requests on the 'reboot' channel
	for {
		select {
		case <-reboot:

			//	Settle a bit before rebooting...
			time.Sleep(3 * time.Second)
			log.Println("[INFO] Rebooting...")
		case <-ctx.Done():
			log.Println("[INFO] Stopping the reboot helper")
			return
		}
	}
}
