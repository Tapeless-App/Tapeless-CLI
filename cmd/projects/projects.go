package projects

import (
	"github.com/spf13/cobra"
	"tapeless.app/tapeless-cli/cmd"
)

func init() {
	cmd.RootCmd.AddCommand(projectsCmd)
}

var (
	projectIdFlag int
	projectsCmd   = &cobra.Command{
		Use:   "projects",
		Short: "Manage your Tapeless projects",
	}
)
