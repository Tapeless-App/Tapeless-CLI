package projects

import (
	"github.com/spf13/cobra"
	"tapeless.app/tapeless-cli/cmd"
	authService "tapeless.app/tapeless-cli/services/auth"
)

func init() {
	cmd.RootCmd.AddCommand(projectsCmd)
}

var (
	projectIdFlag int
	projectsCmd   = &cobra.Command{
		Use:   "projects",
		Short: "Manage your Tapeless projects",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			authService.EnsureValidSession()
		},
	}
)
