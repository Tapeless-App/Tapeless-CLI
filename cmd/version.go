package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"tapeless.app/tapeless-cli/env"
	versionService "tapeless.app/tapeless-cli/services/version"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current Tapeless CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(env.Version)
		versionService.CheckLatestVersion()
	},
}
