package cmd

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

func collectionloop(ctx context.Context) {
	log.Println("[INFO] Running on Windows, so this will be boring...")

	//	Connect to the datastores:
	log.Printf("[INFO] Config database: %s\n", viper.GetString("datastore.config"))
	log.Printf("[INFO] Activities database: %s\n", viper.GetString("datastore.activity"))

	hostname, _ := os.Hostname()
	log.Printf("[INFO] Using hostname %v...", hostname)

	//	Loop and respond to channels:
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(1 * time.Second):
			//	Do nothing?
		}
	}
}
