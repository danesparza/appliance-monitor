package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	buildVersion = "Unknown"
	commitId     string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the version information",
	Long:  `Shows version information and exits`,
	Run: func(cmd *cobra.Command, args []string) {
		//	Show the version number
		fmt.Printf("\nAppliance monitor version %s", buildVersion)

		//	Show the commitid if available:
		if commitId != "" {
			fmt.Printf(" (%s)", commitId[:7])
		}

		//	Trailing space and newline
		fmt.Println(" ")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
