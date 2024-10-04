package repos

import (
	"github.com/spf13/cobra"
	"tapeless.app/tapeless-cli/cmd"
	authService "tapeless.app/tapeless-cli/services/auth"
)

var (
	projectIdFlag int
	repoCmd       = &cobra.Command{
		Use:   "repos",
		Short: "Manage your repositories that are synced with Tapeless",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			authService.EnsureValidSession()
		},
	}
)

func init() {
	cmd.RootCmd.AddCommand(repoCmd)
}
