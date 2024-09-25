package repo

import (
	"fmt"

	"github.com/spf13/cobra"
	reposService "tapeless.app/tapeless-cli/services/repos"
)

func init() {
	repoCmd.AddCommand(listCmd)
	listCmd.Flags().IntVarP(&projectIdFlag, "project-id", "p", -1, "The project ID to list repositories for")
}

var (
	listCmd = &cobra.Command{
		Use:     "list",
		Short:   "List all repositories",
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, args []string) {

			// Sync repositories first
			reposData, err := reposService.GetPersistedRepositories()

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			hasRepos := false

			// Print repositories
			for _, repo := range reposData {

				if projectIdFlag != -1 && projectIdFlag != repo.ProjectId {
					continue
				}

				hasRepos = true

				fmt.Printf("====== Repository: %s ======\n", repo.Name)
				fmt.Println("Path:", repo.Path)
				fmt.Println("Latest Sync:", repo.LatestSync)
				fmt.Println("Project ID:", repo.ProjectId)
				fmt.Println("Git Config ID:", repo.GitConfigId)
				fmt.Println()
			}

			if !hasRepos {
				fmt.Println("No repositories found for project with id:", projectIdFlag)
			}

		},
	}
)
