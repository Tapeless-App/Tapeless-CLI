package repos

import (
	"fmt"

	"github.com/spf13/cobra"
	projectsService "tapeless.app/tapeless-cli/services/projects"
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

			// Get a list of current projects
			projects, err := projectsService.FetchProjects()

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			// Sync repositories first
			reposData, err := reposService.FetchAndUpdateRepositories(projects)

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

				project, err := projectsService.FilterProjectsById(repo.ProjectId, &projects)

				fmt.Printf("====== Repository: %s ======\n", repo.Name)
				fmt.Println("Path:", repo.Path)
				fmt.Println("Latest Sync:", repo.LatestSync)
				fmt.Println("Git Config ID:", repo.GitConfigId)
				if err != nil {
					fmt.Println("Error fetching project with projectId", repo.ProjectId, err)
				} else {
					fmt.Println("Project Name:", project.Name)
				}
				fmt.Println("Project ID:", repo.ProjectId)

				fmt.Println()
			}

			if !hasRepos {
				fmt.Println("No repositories found")
			}

		},
	}
)
