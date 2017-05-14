package cmd

import (
	"log"
	"net/http"
	"strings"

	"github.com/danesparza/appliance-monitor/api"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	serverInterface   string
	serverPort        int
	serverUIDirectory string
	allowedOrigins    string

	maxPoints             = 120
	applianceRunThreshold = 8
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the appliance monitor server",
	Long: `Appliance monitor provides its own webserver which can serve both the 
API and the UI for the app in addition to monitoring the system sensors
and providing notification.`,
	Run: serve,
}

func serve(cmd *cobra.Command, args []string) {

	//	If we have a config file, report it:
	if viper.ConfigFileUsed() != "" {
		log.Println("[INFO] Using config file:", viper.ConfigFileUsed())
	}

	//	Create a router and setup our REST endpoints...
	var Router = mux.NewRouter()

	//	Setup our routes
	// Router.HandleFunc("/", api.ShowUI)
	Router.HandleFunc("/activity/get", api.GetActivity)

	//	Websocket connections
	// Router.Handle("/ws", api.WsHandler{H: api.WsHub})

	//	If we don't have a UI directory specified...
	/*
		if viper.GetString("server.ui-dir") == "" {
			//	Use the static assets file generated with
			//	https://github.com/elazarl/go-bindata-assetfs using the centralconfig-ui from
			//	https://github.com/danesparza/centralconfig-ui.
			//
			//	To generate this file, place the 'ui'
			//	directory under the main centralconfig directory and run the commands:
			//	go-bindata-assetfs.exe -pkg cmd ./ui/...
			//	mv bindata_assetfs.go cmd
			//	go install ./...
			Router.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(assetFS())))
		} else {
			//	Use the supplied directory:
			log.Printf("[INFO] Using UI directory: %s\n", viper.GetString("server.ui-dir"))
			Router.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir(viper.GetString("server.ui-dir")))))
		}
	*/

	//	Setup the CORS options:
	log.Printf("[INFO] Allowed CORS origins: %s\n", viper.GetString("server.allowed-origins"))
	c := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(allowedOrigins, ","),
		AllowCredentials: true,
	})
	corsHandler := c.Handler(Router)

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
	startCmd.Flags().StringVarP(&allowedOrigins, "allowed-origins", "o", "", "comma seperated list of allowed CORS origins")

	//	Bind config flags for optional config file override:
	viper.BindPFlag("server.port", startCmd.Flags().Lookup("port"))
	viper.BindPFlag("server.bind", startCmd.Flags().Lookup("bind"))
	viper.BindPFlag("server.ui-dir", startCmd.Flags().Lookup("ui-dir"))
	viper.BindPFlag("server.allowed-origins", startCmd.Flags().Lookup("allowed-origins"))

}
