package cmd

import (
	"context"
	"log"

	"github.com/danesparza/appliance-monitor/data"
	"github.com/grandcat/zeroconf"
	"github.com/spf13/viper"
)

func zeroconfserver(ctx context.Context) {
	log.Println("[INFO] Starting the zeroconf service...")

	//	Get a reference to the config database
	configDB := data.ConfigDB{
		Database: viper.GetString("datastore.config")}

	//	Get the configured appliance name
	appName, err := configDB.Get("name")
	if err != nil {
		log.Printf("[ERROR] Problem getting appliance name: %v", err)
		return
	}

	//	Create the zeroconf server
	server, err := zeroconf.Register(appName.Value, "appliance.monitor", "local.", 3030, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		log.Printf("[ERROR] Problem starting zeroconf server: %v", err)
		return
	}
	defer server.Shutdown()

	//	Loop and respond to channels:
	for {
		select {
		case <-ctx.Done():
			log.Println("[INFO] Stopping the zeroconf service")
			return
		}
	}

}
