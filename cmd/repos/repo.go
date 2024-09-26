package repos

import (
	"github.com/spf13/cobra"
	"tapeless.app/tapeless-cli/cmd"
)

var (
	projectIdFlag int
	repoCmd       = &cobra.Command{
		Use:   "repos",
		Short: "Manage your repositories that are synced with Tapeless",
	}
)

func init() {
	cmd.RootCmd.AddCommand(repoCmd)
}
