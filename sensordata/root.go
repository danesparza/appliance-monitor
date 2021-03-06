package sensordata

import "github.com/danesparza/appliance-monitor/api"

var (
	// WsHub is the websocket hub for sensor data updates
	WsHub = api.NewHub()

	maxPoints             = 120
	applianceRunThreshold = float64(8)
)
