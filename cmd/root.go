package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	// ProblemWithConfigFile indicates whether or not there was a problem
	// loading the config
	ProblemWithConfigFile bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "appliance-monitor",
	Short: "A service to monitor appliances and provide notifications",
	Long: `Appliance monitor is a REST service to monitor appliances 
using vibration and temperature sensors.  It provides notification 
services, the ability to get recent activity, and wifi configuration
services.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.appliance-monitor.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	//	Set our defaults
	viper.SetDefault("server.port", "3000")
	viper.SetDefault("server.bind", "")
	viper.SetDefault("server.allowed-origins", "*")
	viper.SetDefault("datastore.config", "config.db")
	viper.SetDefault("datastore.activity", "activity.db")

	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.AddConfigPath(".")      // also look in the working directory
	viper.AutomaticEnv()          // read in environment variables that match

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	// If a config file is found, read it in
	// otherwise, make note that there was a problem
	if err := viper.ReadInConfig(); err != nil {
		ProblemWithConfigFile = true
	}
}
