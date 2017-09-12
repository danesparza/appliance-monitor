// +build !linux !arm

package network

import (
	"fmt"
)

// UpdateWifiCredentials changes the hostname for the local machine
func UpdateWifiCredentials(ssid, password string) error {
	return fmt.Errorf("Not running on Linux/ARM, so wifi can't be changed")
}
