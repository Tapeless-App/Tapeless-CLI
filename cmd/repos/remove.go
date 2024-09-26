package repos

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"tapeless.app/tapeless-cli/prompts"
	projectsService "tapeless.app/tapeless-cli/services/projects"
	reposService "tapeless.app/tapeless-cli/services/repos"
)

func init() {
	repoCmd.AddCommand(removeRepoCmd)
	removeRepoCmd.Flags().IntVarP(&projectIdFlag, "project-id", "p", -1, "Project ID from which to remove this repository")
}

var (
	removeRepoCmd = &cobra.Command{
		Use:     "remove",
		Short:   "Remove a repository from your Tapeless project",
		Aliases: []string{"rm"},
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			path := ""

			if len(args) > 0 {

				absPath, err := filepath.Abs(args[0])

				if err != nil {
					fmt.Println("Error getting absolute path to specific path:", err)
					return
				}

				path = absPath

			} else {

				wd, err := os.Getwd()

				if err != nil {
					fmt.Println("No path provided, error getting working directory:", err)
					return
				}

				fmt.Println("No path provided, using working directory:", wd)
				path = wd
			}

			repositories, err := reposService.GetPersistedRepositories()

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			matchingEntries := make([]reposService.Repository, 0)
			entryToRemove := reposService.Repository{}

			for _, repo := range repositories {
				if repo.Path == path {
					// If a projectId flag is provided, only consider repositories with that project ID
					if projectIdFlag != -1 && repo.ProjectId != projectIdFlag {
						continue
					}
					matchingEntries = append(matchingEntries, repo)
				}
			}

			if len(matchingEntries) == 0 {
				fmt.Println("No repositories found with path:", path)
				return
			}

			if len(matchingEntries) == 1 {
				entryToRemove = matchingEntries[0]
			} else {
				if projectIdFlag != -1 {
					fmt.Println("Multiple repositories found with path:", path, "and project ID:", projectIdFlag)
					fmt.Println("This is an unexpected state - deleting the first one found")
					entryToRemove = matchingEntries[0]
				} else {
					fmt.Println("Multiple repositories found with path:", path)
					projects, err := projectsService.SyncProjects()

					if err != nil {
						fmt.Println("Error reading projects:", err)
						return
					}

					matchingProjects := make(map[int]projectsService.ProjectData)

					for _, entry := range matchingEntries {
						matchingProjects[entry.ProjectId] = projects[entry.ProjectId]
					}

					projectId, err := prompts.GetProjectIdPrompt("Select the project to remove the repository from", projectIdFlag, matchingProjects)

					if err != nil {
						fmt.Println("Project selection cancelled")
						return
					}

					for _, entry := range matchingEntries {
						if entry.ProjectId == projectId {
							entryToRemove = entry
							break
						}
					}
				}
			}

			if entryToRemove == (reposService.Repository{}) {
				fmt.Println("No repository found to remove")
				return
			} else {
				fmt.Println("Removing repository:", entryToRemove.Name)
				err = reposService.DeleteGitConfig(entryToRemove.ProjectId, entryToRemove.GitConfigId)
			}

			if err != nil {
				fmt.Println("Error removing repository from project:", err)
				return
			}

			fmt.Println("Repository removed from project successfully - updating local configuration")

			// Remove the repository from the local configuration
			for i, repo := range repositories {
				if repo.GitConfigId == entryToRemove.GitConfigId {
					// Overwrite the repositories slice with the entry removed
					repositories = append(repositories[:i], repositories[i+1:]...)
					break
				}
			}

			viper.Set("repositories", repositories)
			viper.WriteConfig()

		},
	}
)
