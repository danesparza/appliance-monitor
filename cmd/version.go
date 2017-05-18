package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// BuildVersion contains the version information for the app
	BuildVersion = "Unknown"

	// CommitID is the git commitId for the app.  It's filled in as
	// part of the automated build
	CommitID string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the version information",
	Long:  `Shows version information and exits`,
	Run: func(cmd *cobra.Command, args []string) {
		//	Show the version number
		fmt.Printf("\nAppliance monitor version %s", BuildVersion)

		//	Show the commitid if available:
		if CommitID != "" {
			fmt.Printf(" (%s)", CommitID[:7])
		}

		//	Trailing space and newline
		fmt.Println(" ")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
