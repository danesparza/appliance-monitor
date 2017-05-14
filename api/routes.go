package api

import (
	"encoding/json"
	"net/http"

	"github.com/danesparza/appliance-monitor/data"
)

// GetActivity gets activity for the appliance for a given time range
func GetActivity(rw http.ResponseWriter, req *http.Request) {
	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Decode the request:
	request := data.ActivityRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Get the current datastore:
	//	ds := datastores.GetConfigDatastore()

	//	Send the request to the datastore and get a response:
	/*
		response, err := ds.Get(request)
		if err != nil {
			sendErrorResponse(rw, err, http.StatusInternalServerError)
			return
		}
	*/

	//	If we found an item, return it (otherwise, return an empty item):
	sendDataResponse(rw, "Activity found")
}

//	Used to send back an error:
func sendErrorResponse(rw http.ResponseWriter, err error, code int) {
	//	Our return value
	response := data.ActivityResponse{
		Status:  code,
		Message: "Error: " + err.Error()}

	//	Serialize to JSON & return the response:
	rw.WriteHeader(code)
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

//	Used to send back a response with data
func sendDataResponse(rw http.ResponseWriter, message string) {
	//	Our return value
	response := data.ActivityResponse{
		Status:  http.StatusOK,
		Message: message,
		ID:      3} // TODO: Pass data back in a more flexible way

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
