package api

import (
	"encoding/json"
	"net/http"

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
