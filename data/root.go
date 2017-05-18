package data

import "time"

// CurrentState describes the current running state of the application
type CurrentState struct {
	ServerStartTime    time.Time `json:"starttime"`
	ApplicationVersion string    `json:"appversion"`
	DeviceRunning      bool      `json:"devicerunning"`
}
