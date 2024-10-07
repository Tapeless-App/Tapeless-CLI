package projects

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"tapeless.app/tapeless-cli/prompts"
	projectService "tapeless.app/tapeless-cli/services/projects"
)

func init() {
	projectsCmd.AddCommand(removeProjectCmd)
	removeProjectCmd.Flags().IntVarP(&projectIdFlag, "project-id", "p", -1, "The project ID to list repositories for")
}

var removeProjectCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove a project",
	Aliases: []string{"rm", "delete"},
	Run: func(cmd *cobra.Command, args []string) {

		projects, err := projectService.FetchProjects()

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		project, err := prompts.GetProjectIdPrompt("Select the project you wish to delete", projectIdFlag, projects)

		if err != nil {
			fmt.Println("Project selection cancelled")
			return
		}

		confirmationPrompt := promptui.Prompt{
			Label:     fmt.Sprintf("Are you sure you want to remove project '%s'?", project.Name),
			IsConfirm: true,
		}

		_, err = confirmationPrompt.Run()

		if err != nil {
			fmt.Println("Project removal cancelled")
			return
		}

		err = projectService.DeleteProject(project.Id)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Project removed successfully")

	},
}
