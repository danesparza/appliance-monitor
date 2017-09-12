package network

import (
	"bytes"
	"html/template"
)

var wifiSupplicant = []byte(`country=US
	
network={
	ssid="{{.SSID}}"
	psk="{{.Passphrase}}"
	key_mgmt=WPA-PSK
}

`)

// WifiInfo describes the information needed to configure the wifi client
type WifiInfo struct {
	SSID       string
	Passphrase string
}

// FormatWifiCredentials formats the given credentials as a new wifi supplicant config file
func FormatWifiCredentials(ssid, password string) (string, error) {

	credentials := WifiInfo{ssid, password}
	tmpl, err := template.New("wifisupp").Parse(string(wifiSupplicant))
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, credentials)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil
}
