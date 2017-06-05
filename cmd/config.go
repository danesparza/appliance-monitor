package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	jsonConfig bool
	yamlConfig bool
)

var yamlDefault = []byte(`
server:
  bind: "127.0.0.1"
  port: 3030
  allowed-origins: "*"
settings:
  name: "appliance-monitor"
  monitorwindow: 120
`)

var jsonDefault = []byte(`{
	"server": {
		"bind": "127.0.0.1",
		"port": 3030,
		"allowed-origins": "*"
	},
	"settings": {
		"name": "appliance-monitor",
		"monitorwindow": 120
	}
}`)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Prints default server configuration files",
	Long: `Use this to create a default configuration file for the appliance monitor  

Example: 

appliance-monitor config > config.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		if jsonConfig {
			fmt.Printf("%s", jsonDefault)
		} else if yamlConfig {
			fmt.Printf("%s", yamlDefault)
		}
	},
}

func init() {
	RootCmd.AddCommand(configCmd)

	configCmd.Flags().BoolVarP(&jsonConfig, "json", "j", false, "Create a JSON configuration file")
	configCmd.Flags().BoolVarP(&yamlConfig, "yaml", "y", true, "Create a YAML configuration file")

}
