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
}

// GetCurrentState gets the current running state of the application
func GetCurrentState(rw http.ResponseWriter, req *http.Request) {
	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Find out if the device is currently running:
	activityDB := data.ActivityDB{
		Database: viper.GetString("datastore.activity")}
	response, _ := activityDB.GetLatestActivity()

	//	Get the current datastore:
	currentState := CurrentState{
		ServerStartTime:    ApplicationStartTime,
		DeviceRunning:      response.Type == data.ApplianceRunning,
		ApplicationVersion: BuildVersion,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(currentState)
}
