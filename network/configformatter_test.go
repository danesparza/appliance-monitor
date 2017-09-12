package network_test

import (
	"testing"

	"github.com/danesparza/appliance-monitor/network"
)

//	The formatting utility should format the wifi supplicant file
func TestFormatWifiCredentials_WithValidParams_ShouldReturnProperFormat(t *testing.T) {

	//	Arrange
	ssid := "testssid"
	passphrase := "testpassphrase"
	expectedconfig := `country=US
	
network={
	ssid="testssid"
	psk="testpassphrase"
	key_mgmt=WPA-PSK
}

`

	//	Act
	retval, err := network.FormatWifiCredentials(ssid, passphrase)

	//	Assert
	if err != nil {
		t.Errorf("An error occured while formatting the wifi config: %v", err)
	}

	if retval != expectedconfig {
		t.Errorf("Configuration doesn't match what we expect.  Here's what we got:\n%v", retval)
	}
}
