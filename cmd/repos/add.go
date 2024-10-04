package repos

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"tapeless.app/tapeless-cli/prompts"
	projectService "tapeless.app/tapeless-cli/services/projects"
	repoService "tapeless.app/tapeless-cli/services/repos"
)

func init() {
	repoCmd.AddCommand(addCmd)
}

var (
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a new repository to be synced with Tapeless",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			path := ""

			if len(args) > 0 {

				path = args[0]

			} else {

				wd, err := os.Getwd()

				if err != nil {
					fmt.Println("No path provided, error getting working directory:", err)
					return
				}

				pathPrompt := promptui.Prompt{
					Label:     "Enter the path to the repository",
					Default:   wd,
					AllowEdit: true,
					Validate: func(input string) error {
						if _, err := os.Stat(input); os.IsNotExist(err) {
							return fmt.Errorf("path '%s' does not exist", input)
						}
						return nil
					},
				}

				path, err = pathPrompt.Run()

				if err != nil {
					fmt.Println("Path selection cancelled")
					return
				}
			}

			projects, err := projectService.FetchProjects()

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			projectId, err := prompts.GetProjectIdPrompt("Select the project this repository belongs to", projectIdFlag, projects)

			if err != nil {
				fmt.Println("Project selection cancelled")
				return
			}

			repositories := make([]repoService.Repository, 0)

			viper.UnmarshalKey("repositories", &repositories)

			localGitConfig, err := repoService.GetLocalRepositoryConfig(path)

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			project, err := projectService.FilterProjectsById(projectId, &projects)

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			projectName := project.Name
			localGitConfig.Name = strings.Join([]string{projectName, localGitConfig.Name}, "/")

			fmt.Println("Current path:", localGitConfig.Path)
			fmt.Println("Email Name:", localGitConfig.AuthorEmail)
			fmt.Println("Repo Name:", localGitConfig.Name)
			fmt.Println("Origin:", localGitConfig.OriginUrl)
			fmt.Println("Git Config Name:", localGitConfig.Name)

			existingConfigIndex := slices.IndexFunc(repositories, func(r repoService.Repository) bool {
				return r.ProjectId == projectId && r.Path == localGitConfig.Path
			})

			if existingConfigIndex != -1 {
				fmt.Println("Repository already exists in the project!")
				return
			}

			gitConfigResponse, err := repoService.CreateGitConfig(projectId, localGitConfig)

			if err != nil {
				fmt.Println("Error creating git config:", err)
				return
			}

			repositories = append(repositories, repoService.Repository{
				Name:        gitConfigResponse.Name,
				Path:        localGitConfig.Path,
				LatestSync:  "",
				ProjectId:   gitConfigResponse.ProjectId,
				GitConfigId: gitConfigResponse.Id,
				AuthorEmail: localGitConfig.AuthorEmail,
				OriginUrl:   localGitConfig.OriginUrl,
			})

			viper.Set("repositories", repositories)
			viper.WriteConfig()

		},
	}
)
