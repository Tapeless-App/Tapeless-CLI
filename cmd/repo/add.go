package repo

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	projectService "tapeless.app/tapeless-cli/services/projects"
	repoService "tapeless.app/tapeless-cli/services/repos"
	"tapeless.app/tapeless-cli/util"
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

				fmt.Println("No path provided, using working directory:", wd)
				path = wd
			}

			projects, err := projectService.SyncProjects()

			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			projectId, err := getProjectId(projectIdFlag, projects)

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

			projectName := projects[projectId].Name
			gitConfigName := strings.Join([]string{projectName, localGitConfig.Name}, "/")

			fmt.Println("Current path:", localGitConfig.Path)
			fmt.Println("Email Name:", localGitConfig.AuthorEmail)
			fmt.Println("Repo Name:", localGitConfig.Name)
			fmt.Println("Origin:", localGitConfig.OriginUrl)
			fmt.Println("Git Config Name:", gitConfigName)

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
			})

			viper.Set("repositories", repositories)
			viper.WriteConfig()

		},
	}
)

/**
 * Get the project ID for the repository
 * Will use the flag if it is set, otherwise prompt the user with a list of projects
 */
func getProjectId(projectIdFlag int, projects map[int]projectService.ProjectData) (int, error) {

	if projectIdFlag != -1 {
		return projectIdFlag, nil
	}

	items := []projectService.ProjectData{}

	for _, project := range projects {

		items = append(items, projectService.ProjectData{
			Id:           project.Id,
			Name:         project.Name,
			LastSync:     util.FormatDate(project.LastSync),
			ProjectStart: util.FormatDate(project.ProjectStart),
			ProjectEnd:   util.FormatDate(project.ProjectEnd),
		})
	}

	templates := &promptui.SelectTemplates{
		Label:    `{{ . }}:`,
		Active:   "> {{ .Name | cyan }} (id: {{ .Id }})",
		Inactive: "  {{ .Name | cyan }} (id: {{ .Id }})",
		Selected: "{{ .Name }}",
		Details: `
--------- Project: {{ .Name }} ----------

{{ "Id:" | faint }}	{{ .Id }}
{{ "Project Start:" | faint }}	{{ .ProjectStart }}
{{ "Project End:" | faint }}	{{ .ProjectEnd }}
{{ "Last Sync:" | faint }}	{{ .LastSync }}
 
`,
	}

	prompt := promptui.Select{
		Templates: templates,
		Label:     "Select the project this repository belongs to",
		Items:     items,
		Size:      len(items),
	}

	index, _, err := prompt.Run()

	if err != nil {
		return -1, err
	}

	return items[index].Id, nil
}
