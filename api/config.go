package api

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/danesparza/appliance-monitor/data"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// GetAllConfig gets all configuration and returns it in JSON format.
func GetAllConfig(rw http.ResponseWriter, req *http.Request) {

	//	Connect to the datastore:
	configDB := data.ConfigDB{
		Database: viper.GetString("datastore.config")}

	response, err := configDB.GetAll()

	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// GetConfigItem gets a single configuration and returns it in JSON format.
// If the item can't be found, returns an empty config item
func GetConfigItem(rw http.ResponseWriter, req *http.Request) {

	//	Connect to the datastore:
	configDB := data.ConfigDB{
		Database: viper.GetString("datastore.config")}

	//	Get the config name from the request:
	configName := mux.Vars(req)["name"]

	//	See if we can find a config item with that name:
	response, err := configDB.Get(configName)

	if err != nil {
		sendErrorResponse(rw, err, 500)
		return
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// RemoveConfigItem removes a single config item
func RemoveConfigItem(rw http.ResponseWriter, req *http.Request) {

	//	Get the config datastore:
	configDB := data.ConfigDB{
		Database: viper.GetString("datastore.config")}

	//	Get the config name from the request:
	configName := mux.Vars(req)["name"]

	//	Send the request to the datastore and get a response:
	err := configDB.Remove(configName)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(fmt.Sprintf("Removed %s", configName))
}

// SetConfigItem adds or updates a single config item and returns the new item in JSON format
func SetConfigItem(rw http.ResponseWriter, req *http.Request) {
	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Decode the request if it was a POST:
	request := data.ConfigItem{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Get the config datastore:
	configDB := data.ConfigDB{
		Database: viper.GetString("datastore.config")}

	//	Send the request to the datastore and get a response:
	response, err := configDB.Set(request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// SetAllConfigItems adds or updates multiple config items and returns all config items in JSON format
func SetAllConfigItems(rw http.ResponseWriter, req *http.Request) {
	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Decode the request if it was a POST:
	request := []data.ConfigItem{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Get the config datastore:
	configDB := data.ConfigDB{
		Database: viper.GetString("datastore.config")}

	//	Send each request to the datastore and get a response:
	for c := 0; c < len(request); c++ {
		//	Set the config item:
		_, err := configDB.Set(request[c])
		if err != nil {
			sendErrorResponse(rw, err, http.StatusInternalServerError)
			return
		}
	}

	//	Get the (updated) full list of config items:
	response, err := configDB.GetAll()
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
