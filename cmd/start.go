package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/danesparza/appliance-monitor/api"
	"github.com/danesparza/appliance-monitor/sensordata"
	"github.com/danesparza/appliance-monitor/zeroconf"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	serverInterface   string
	serverPort        int
	serverUIDirectory string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the appliance monitor server",
	Long: `Appliance monitor provides its own webserver which can serve both the 
API and the UI for the app in addition to monitoring the system sensors
and providing notification.`,
	Run: start,
}

func start(cmd *cobra.Command, args []string) {

	//	If we have a config file, report it:
	if viper.ConfigFileUsed() != "" {
		log.Println("[INFO] Using config file:", viper.ConfigFileUsed())
	}

	//	Create a router and setup our REST endpoints...
	Router := mux.NewRouter()

	//	Setup our routes
	Router.HandleFunc("/", api.ShowUI)

	//	Activities
	Router.HandleFunc("/activity", api.GetAllActivity).Methods("GET")
	Router.HandleFunc("/activity", api.GetActivityInRange).Methods("POST")

	//	Config
	Router.HandleFunc("/config", api.GetAllConfig).Methods("GET")
	Router.HandleFunc("/config/{name}", api.GetConfigItem).Methods("GET")
	Router.HandleFunc("/config", api.SetConfigItem).Methods("POST")
	Router.HandleFunc("/config/{name}", api.RemoveConfigItem).Methods("DELETE")

	//	System information
	Router.HandleFunc("/system/state", api.GetCurrentState).Methods("GET")
	Router.HandleFunc("/system/wifi", nil)

	//	System resets
	Router.HandleFunc("/reset/network", nil)
	Router.HandleFunc("/reset/config", nil)

	//	Websocket connections
	Router.Handle("/ws", api.WsHandler{H: sensordata.WsHub})

	//	Trap program exit appropriately
	ctx, cancel := context.WithCancel(context.Background())
	sigch := make(chan os.Signal, 2)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	go handleSignals(ctx, sigch, cancel)

	//	Start the collection process
	go sensordata.CollectAndProcess(ctx)

	//	Start the zeroconf server
	go zeroconf.Serve(ctx)

	//	If we don't have a UI directory specified...
	if viper.GetString("server.ui-dir") == "" {
		//	Use the static assets file generated with
		//	https://github.com/elazarl/go-bindata-assetfs using the application-monitor-ui from
		//	https://github.com/danesparza/application-monitor-ui.
		//
		//	To generate this file, place the 'ui'
		//	directory under the main application-monitor-ui directory and run the commands:
		//	go-bindata-assetfs -pkg cmd build/...
		//	Move bindata_assetfs.go to the application-monitor cmd directory
		//	go install ./...
		Router.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(assetFS())))
	} else {
		//	Use the supplied directory:
		log.Printf("[INFO] Using UI directory: %s\n", viper.GetString("server.ui-dir"))
		Router.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir(viper.GetString("server.ui-dir")))))
	}

	//	Setup the CORS options:
	log.Printf("[INFO] Allowed CORS origins: %s\n", viper.GetString("server.allowed-origins"))

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(viper.GetString("server.allowed-origins"), ","),
		AllowCredentials: true,
	}).Handler(Router)

	//	Format the bound interface:
	formattedInterface := viper.GetString("server.bind")
	if formattedInterface == "" {
		formattedInterface = "127.0.0.1"
	}

	//	If we have an SSL cert specified, use it:
	if viper.GetString("server.sslcert") != "" {
		log.Printf("[INFO] Using SSL cert: %s\n", viper.GetString("server.sslcert"))
		log.Printf("[INFO] Using SSL key: %s\n", viper.GetString("server.sslkey"))
		log.Printf("[INFO] Starting HTTPS server: https://%s:%s\n", formattedInterface, viper.GetString("server.port"))

		log.Printf("[ERROR] %v\n", http.ListenAndServeTLS(viper.GetString("server.bind")+":"+viper.GetString("server.port"), viper.GetString("server.sslcert"), viper.GetString("server.sslkey"), corsHandler))
	} else {
		log.Printf("[INFO] Starting HTTP server: http://%s:%s\n", formattedInterface, viper.GetString("server.port"))
		log.Printf("[ERROR] %v\n", http.ListenAndServe(viper.GetString("server.bind")+":"+viper.GetString("server.port"), corsHandler))
	}
}

func init() {
	RootCmd.AddCommand(startCmd)

	//	Setup our flags
	startCmd.Flags().IntVarP(&serverPort, "port", "p", 1313, "port on which the server will listen")
	startCmd.Flags().StringVarP(&serverInterface, "bind", "i", "", "interface to which the server will bind")
	startCmd.Flags().StringVarP(&serverUIDirectory, "ui-dir", "u", "", "directory for the UI")

	//	Bind config flags for optional config file override:
	viper.BindPFlag("server.port", startCmd.Flags().Lookup("port"))
	viper.BindPFlag("server.bind", startCmd.Flags().Lookup("bind"))
	viper.BindPFlag("server.ui-dir", startCmd.Flags().Lookup("ui-dir"))

	//	Set the build version information:
	api.BuildVersion = BuildVersion
	api.CommitID = CommitID

}

func handleSignals(ctx context.Context, sigch <-chan os.Signal, cancel context.CancelFunc) {
	select {
	case <-ctx.Done():
	case sig := <-sigch:
		switch sig {
		case os.Interrupt:
			log.Println("[INFO] SIGINT")
		case syscall.SIGTERM:
			log.Println("[INFO] SIGTERM")
		}
		log.Println("[INFO] Shutting down ...")
		cancel()
		os.Exit(0)
	}
}
