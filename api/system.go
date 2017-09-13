package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/danesparza/appliance-monitor/data"
	"github.com/danesparza/appliance-monitor/network"
	"github.com/spf13/viper"
)

// CurrentState describes the current running state of the application
type CurrentState struct {
	ServerStartTime    time.Time `json:"starttime"`
	ApplicationVersion string    `json:"appversion"`
	DeviceRunning      bool      `json:"devicerunning"`
	DeviceID           string    `json:"deviceId"`
}

// WifiUpdateRequest describes the request to update the Wifi credentials
type WifiUpdateRequest struct {
	SSID       string `json:"ssid"`
	Passphrase string `json:"passphrase"`
}

type WifiUpdateResponse struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
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

// UpdateWifi updates the wifi credentials for the machine
func UpdateWifi(rw http.ResponseWriter, req *http.Request) {

	//	Decode the request if it was a POST:
	request := WifiUpdateRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Send the request to the wifi helper:
	response := WifiUpdateResponse{Status: 200, Description: "Successful.  Rebooting..."}
	err = network.UpdateWifiCredentials(request.SSID, request.Passphrase)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)

	//	Reboot the machine
	go network.RebootMachine()
}
