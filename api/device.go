package api

import "time"

// Device represents a device in the system.
type Device struct {
	// ID is the unique device identifier
	ID string `json:"id"`

	// Type is the type of device (amon / hs110 / etc)
	Type string `json:"type"`

	// Name is the name given to the device by the customer
	Name string `json:"name"`

	// MinimumMonitorTime is the amount of time required to monitor the device before
	// it latches to the 'on' or 'off' state
	MinimumMonitorTime time.Duration `json:"minimum_monitor_time"`

	// Threshold is the value that must be crossed (after the MinimumMonitorTime time)
	// before the device latches to the 'on' or 'off' state
	Threshold int `json:"monitor_threshold"`

	// IPAddress is the network address for the device
	IPAddress string `json:"ipaddress"`

	// Running indicates whether the device is currently running (operating) or not
	Running bool `json:"running"`
}
