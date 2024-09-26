package projects

import (
	"fmt"

	"github.com/spf13/cobra"
	projectsService "tapeless.app/tapeless-cli/services/projects"
)

func init() {
	projectsCmd.AddCommand(listPorjectsCmd)
}

var listPorjectsCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all projects",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {

		// Sync projects first
		projectsData, err := projectsService.SyncProjects()

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Print projects
		for i := len(projectsData) - 1; i >= 0; i-- {
			project := &projectsData[i]
			fmt.Printf("====== Project: %s ======\n", project.Name)
			fmt.Println("ID:", project.Id)
			fmt.Println("Project Start:", project.ProjectStart)
			fmt.Println("Project End:", project.ProjectEnd)
			fmt.Println("Last Sync:", project.LastSync)
			fmt.Println()
		}

	},
}
