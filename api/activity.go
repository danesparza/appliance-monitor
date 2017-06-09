package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/araddon/dateparse"
	"github.com/danesparza/appliance-monitor/data"
	"github.com/spf13/viper"
)

// ActivityRequest represents an API request for activity
type ActivityRequest struct {
	StartTime string `json:"starttime"`
	EndTime   string `json:"endtime"`
}

// GetActivity gets activity for the appliance for a given time range
func GetActivity(rw http.ResponseWriter, req *http.Request) {
	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	The default request:
	starttime := time.Now().Add(-24 * time.Hour)
	endtime := time.Now()

	//	Decode the request if it was a POST:
	if req.Method == "POST" {
		request := ActivityRequest{}
		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			sendErrorResponse(rw, err, http.StatusBadRequest)
			return
		}

		//	Parse the dates from many different formats: https://github.com/araddon/dateparse
		if t, err := dateparse.ParseAny(request.StartTime); err == nil {
			starttime = t
		}

		if t, err := dateparse.ParseAny(request.EndTime); err == nil {
			endtime = t
		}
	}

	//	Get the activity datastore:
	activityDB := data.ActivityDB{
		Database: viper.GetString("datastore.activity")}

	//	Send the request to the datastore and get a response:
	response, err := activityDB.GetRange(starttime, endtime)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
