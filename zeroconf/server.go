package zeroconf

import (
	"context"
	"log"

	"fmt"

	"github.com/danesparza/appliance-monitor/data"
	"github.com/grandcat/zeroconf"
	"github.com/rs/xid"
	"github.com/spf13/viper"
)

// Serve starts the zeroconf service and registers this server
func Serve(ctx context.Context, restart chan bool) {
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

	//	Get our machine id:
	guid := xid.New()
	machineID := guid.Machine()

	//	Create the zeroconf server
	server, err := zeroconf.Register(appName.Value, "_appliance-monitor._tcp", "local.", 3030, []string{
		"txtv=1", fmt.Sprintf("id=%x", machineID)}, nil)

	if err != nil {
		log.Printf("[ERROR] Problem starting zeroconf server: %v", err)
		return
	}
	defer server.Shutdown()

	//	Loop and respond to channels:
	for {
		select {
		case <-restart:

			//	Get the configured appliance name
			appName, err = configDB.Get("name")
			if err != nil {
				log.Printf("[ERROR] Problem refetching appliance name: %v", err)
				return
			}

			log.Printf("[INFO] Restarting service with name: %v", appName.Value)

			//	Shutdown the old server
			server.Shutdown()

			//	Start a new server with the new name
			server, err = zeroconf.Register(appName.Value, "_appliance-monitor._tcp", "local.", 3030, []string{
				"txtv=1", fmt.Sprintf("id=%x", machineID)}, nil)

			if err != nil {
				log.Printf("[ERROR] Problem restarting zeroconf server: %v", err)
				return
			}
		case <-ctx.Done():
			log.Println("[INFO] Stopping the zeroconf service")
			return
		}
	}

}
