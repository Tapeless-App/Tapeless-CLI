package repos

import (
	"fmt"

	"github.com/spf13/cobra"
	reposService "tapeless.app/tapeless-cli/services/repos"
)

func init() {
	repoCmd.AddCommand(syncCmd)
}

var (
	syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Sync all repositories with Tapeless",
		Run: func(cmd *cobra.Command, args []string) {

			// Sync repositories
			repositories, err := reposService.GetPersistedRepositories()

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			if len(repositories) == 0 {
				fmt.Println("No repositories found - add a repository first using \"tapeless repos add\"")
				return
			}

			testRepo := repositories[0]

			fmt.Println("Syncing repository:", testRepo.Name)

			// url := fmt.Sprintf("http://localhost:4000/projects/%d/gitConfigs/%d/commits", testRepo.ProjectId, testRepo.GitConfigId)

		},
	}
)
