// Perhaps better "create" or "new" instead of "add"? What is more idiomatic in conjunction with "repos"?

package projects

import (
	"errors"
	"fmt"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	projectsService "tapeless.app/tapeless-cli/services/projects"
	"tapeless.app/tapeless-cli/util"
)

func init() {

	projectsCmd.AddCommand(addProjectsCmd)

}

var addProjectsCmd = &cobra.Command{
	Use:     "add",
	Short:   "Create a new project",
	Aliases: []string{"new", "create"},
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Creating a new project")

		projectNamePrompt := promptui.Prompt{
			Label: "What is the NAME of the project?",
			Validate: func(input string) error {
				if len(input) < 3 {
					return errors.New("project name must be at least 3 characters long")
				}
				return nil
			},
		}

		projectStartPrompt := promptui.Prompt{
			Label: "When does the project START (format yyyy-mm-dd)?",
			Validate: func(input string) error {
				_, err := time.Parse("2006-01-02", input)

				if err != nil {
					return errors.New("invalid date format - must be yyyy-mm-dd")
				}

				return nil
			},
			Default:   time.Now().Format("2006-01-02"),
			AllowEdit: true,
		}

		projectEndPrompt := promptui.Prompt{
			Label: "Optional: When does the project END (format yyyy-mm-dd)?",
			Validate: func(input string) error {
				if input != "" {
					_, err := time.Parse("2006-01-02", input)

					if err != nil {
						return errors.New("invalid date format - must be yyyy-mm-dd")
					}
				}

				return nil
			},
			Default: "",
		}

		projectName, err := projectNamePrompt.Run()

		if err != nil {
			fmt.Println("Error in setting project name:", err)
			return
		}

		projectStart, err := projectStartPrompt.Run()

		if err != nil {
			fmt.Println("Error in setting project start date:", err)
			return
		}

		projectEnd, err := projectEndPrompt.Run()

		if err != nil {
			fmt.Println("Error in setting project end date:", err)
			return
		}

		err = util.VerifyStartBeforeEnd(projectStart, projectEnd, "2006-01-02")

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		Project, err := projectsService.CreateProject(
			projectsService.ProjectsCreateRequest{
				Name:         projectName,
				ProjectStart: projectStart,
				ProjectEnd:   projectEnd,
			})

		if err != nil {
			fmt.Println("Error creating project:", err)
			return
		}

		fmt.Printf("Project '%s' (id: '%d') created successfully!\n", Project.Name, Project.Id)

		fmt.Println("Syncing projects...")

		_, err = projectsService.FetchProjects()

		if err != nil {
			fmt.Println("Error syncing projects:", err)
			return
		}

		fmt.Println("Projects synced successfully!")

	},
}
