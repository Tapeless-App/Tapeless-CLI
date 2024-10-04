package repos

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
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

			projects, err := projectsService.FetchProjects()

			if err != nil {
				fmt.Println("Error fetching projects:", err)
				return
			}

			repositories, err := reposService.FetchAndUpdateRepositories(projects)

			if err != nil {
				fmt.Println("Error fetching repositories:", err)
				return
			}

			if len(repositories) == 0 {
				fmt.Println("No repositories found")
				return
			}

			// Get the repository to remove
			repoToRemove, err := prompts.GetRepositoryPrompt("Select the repository to remove", repositories, projects)

			if err != nil {
				fmt.Println("Repository removal cancelled", err)
				return
			}

			confirmationPrompt := promptui.Prompt{
				Label:     fmt.Sprintf("Are you sure you want to remove the repository '%s'?", repoToRemove.Name),
				IsConfirm: true,
				Default:   "n",
			}

			_, err = confirmationPrompt.Run()

			if err != nil {
				fmt.Println("Repository removal cancelled")
				return
			}

			reposService.DeleteRepository(*repoToRemove)

			fmt.Println("Repository removed from project successfully - updating local configuration")

			// Remove the repository from the local configuration
			for i, repo := range repositories {
				if repo.GitConfigId == repoToRemove.GitConfigId {
					// Overwrite the repositories slice with the entry removed
					repositories = append(repositories[:i], repositories[i+1:]...)
					break
				}
			}

			err = reposService.PersistRepositories(repositories)

			if err != nil {
				fmt.Println("Error persisting repositories:", err)
				return
			}

		},
	}
)
