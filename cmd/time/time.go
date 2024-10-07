package time

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"tapeless.app/tapeless-cli/cmd"
	"tapeless.app/tapeless-cli/prompts"
	authService "tapeless.app/tapeless-cli/services/auth"
	projectsService "tapeless.app/tapeless-cli/services/projects"
	reposService "tapeless.app/tapeless-cli/services/repos"
)

var (
	dateFlag      string
	projectIdFlag int
	timeCmd       = &cobra.Command{
		Use:   "time",
		Short: "Add and manage time entries",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			authService.EnsureValidSession()
		},
	}
)

func init() {
	cmd.RootCmd.AddCommand(timeCmd)
	timeCmd.PersistentFlags().IntVarP(&projectIdFlag, "project-id", "p", -1, "The project ID for which to manage your time entries")
	timeCmd.PersistentFlags().StringVarP(&dateFlag, "date", "d", "", "The date for which to manage your time entries (format: yyyy-mm-dd)")
}

// Fetches a project based on the following criteria:
// 1. If there is a projectIdFlag != -1, then use that project
// 2. If the working directory is a repo that is added to a project, use that project
// 3. If the above criteria do not match, prompt the user to select a project
func GetProjectBasedOnWorkingDir(backupPromptLabel string, projectIdFlag int) (projectsService.Project, error) {

	projects, err := projectsService.FetchProjects()

	if err != nil {
		return projectsService.Project{}, err
	}

	wd, err := os.Getwd()
	var project projectsService.Project

	if err != nil {
		project, err = prompts.GetProjectIdPrompt(
			"Select a project to create a time entry for",
			-1,
			projects)

		if err != nil {

			return projectsService.Project{}, err
		}
	} else {

		repos, err := reposService.FetchAndUpdateRepositories(projects)

		if err != nil {
			fmt.Println("Error fetching repositories - select project manually")

			project, err = prompts.GetProjectIdPrompt(
				"Select a project to create a time entry for",
				-1,
				projects)

			if err != nil {
				return projectsService.Project{}, err
			}
		} else {
			repo, err := reposService.GetRepositoryByDir(wd, repos)

			if err != nil {
				fmt.Println("No repository setup for current working directory - select project manually")

				project, err = prompts.GetProjectIdPrompt(
					"Select a project to create a time entry for",
					-1,
					projects)

				if err != nil {
					return projectsService.Project{}, err
				}
			} else {
				return prompts.GetProjectIdPromptWithDefault(
					"Select a project to create a time entry for",
					projectIdFlag,
					projects,
					repo.ProjectId)
			}
		}

	}

	return project, nil
}

func GetDate(promptLabel string, dateFlag string) (string, error) {
	if _, err := time.Parse("2006-01-02", dateFlag); err != nil || dateFlag == "" {

		if err != nil && dateFlag != "" {
			println("Invalid date format - must be yyyy-mm-dd")
		}

		datePrompt := promptui.Prompt{
			Label:   promptLabel,
			Default: time.Now().Format("2006-01-02"),
			Validate: func(input string) error {
				_, err := time.Parse("2006-01-02", input)

				if err != nil {
					return errors.New("invalid date format - must be yyyy-mm-dd")
				}

				return nil
			},

			AllowEdit: true,
		}

		return datePrompt.Run()
	} else {
		return dateFlag, nil
	}
}
