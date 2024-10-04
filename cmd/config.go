package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"tapeless.app/tapeless-cli/env"
	projectsService "tapeless.app/tapeless-cli/services/projects"
	reposService "tapeless.app/tapeless-cli/services/repos"
)

func init() {
	RootCmd.AddCommand(configCmd)
}

var (
	configCmd = &cobra.Command{
		Use:     "config",
		Aliases: []string{"debug"},
		Short:   "Get the local Tapeless configuration for debug purposes",
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println("=== Tapeless configuration ===")
			fmt.Println("Version:", env.Version)
			fmt.Println("API URL:", env.ApiURL)
			fmt.Println("Web URL:", env.WebURL)
			fmt.Println("Login Callback Port:", env.LoginCallbackPort)
			fmt.Println("Token:", viper.GetString("token"))
			fmt.Println("")
			projects, err := projectsService.FetchProjects()

			if err != nil {
				fmt.Println("Error fetching projects:", err)
				return
			}

			if len(projects) == 0 {
				fmt.Println("No projects found")
				return
			}

			// Print projects
			for i := len(projects) - 1; i >= 0; i-- {
				project := &projects[i]
				fmt.Printf("====== Project: %s ======\n", project.Name)
				fmt.Println("ID:", project.Id)
				fmt.Println("Project Start:", project.ProjectStart)
				fmt.Println("Project End:", project.ProjectEnd)
				fmt.Println("Last Sync:", project.LastSync)
				fmt.Println()
			}

			repos, err := reposService.FetchAndUpdateRepositories(projects)

			if err != nil {
				fmt.Println("Error fetching repositories:", err)
				return
			}

			if len(repos) == 0 {
				fmt.Println("No repositories found")
				return
			}

			for _, repo := range repos {
				fmt.Printf("====== Repository: %s ======\n", repo.Name)
				fmt.Println("Path:", repo.Path)
				fmt.Println("Latest Sync:", repo.LatestSync)
				fmt.Println("Project ID:", repo.ProjectId)
				fmt.Println("Git Config ID:", repo.GitConfigId)
				fmt.Println()
			}
		},
	}
)
