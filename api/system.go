package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/danesparza/appliance-monitor/data"
	"github.com/spf13/viper"
)

// CurrentState describes the current running state of the application
type CurrentState struct {
	ServerStartTime    time.Time `json:"starttime"`
	ApplicationVersion string    `json:"appversion"`
	DeviceRunning      bool      `json:"devicerunning"`
	DeviceID           string    `json:"deviceId"`
}

// GetCurrentState gets the current running state of the application
func GetCurrentState(rw http.ResponseWriter, req *http.Request) {

	//	Find out if the device is currently running:
	activityDB := data.ActivityDB{
		Database: viper.GetString("datastore.activity")}
	latestActivity, _ := activityDB.GetLatestActivity()

	//	Get config information:
	configDB := data.ConfigDB{
		Database: viper.GetString("datastore.config")}
	deviceID, _ := configDB.Get("deviceID")

	//	Create a CurrentState type:
	currentState := CurrentState{
		ServerStartTime:    ApplicationStartTime,
		DeviceRunning:      latestActivity.Type == data.ApplianceRunning,
		ApplicationVersion: BuildVersion,
		DeviceID:           deviceID.Value,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(currentState)
}
