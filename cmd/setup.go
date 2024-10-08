package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	projectsService "tapeless.app/tapeless-cli/services/projects"
	reposService "tapeless.app/tapeless-cli/services/repos"
	"tapeless.app/tapeless-cli/util"
)

func init() {
	RootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:     "setup",
	Aliases: []string{"status"},
	Short:   "Get your Tapeless setup and configuration status",
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("Note: You can run 'tapeless setup' at any time to get the next steps.")
	},
	Run: func(cmd *cobra.Command, args []string) {

		token := viper.GetString("token")

		fmt.Println("=== Login Status ===")

		if token == "" {
			fmt.Println("Not logged in.")
			fmt.Println("Next step: Run 'tapeless login' login to tapeless.")
			return
		}

		isExpired, err := util.IsJWTExpired(token)

		if err != nil {
			fmt.Println("Error verifying access token:", err)
			fmt.Println("Next step: Try to run 'tapeless login' - if that doesn't work try deleting the config file and starting over.")
			return
		}

		if isExpired {
			fmt.Println("JWT token expired.")
			fmt.Println("Next step: Run 'tapeless login' login to tapeless.")
			return
		}

		fmt.Println("You are logged in.")
		fmt.Println()

		fmt.Println("=== Project Setup ===")

		projects, err := projectsService.FetchProjects()

		if err != nil {
			fmt.Println("Error reading projects:", err)
			return
		}

		if len(projects) == 0 {
			fmt.Println("No projects found.")
			fmt.Println("Next step: Run 'tapeless projects add' to add a project.")
			return
		}

		fmt.Println("Found", len(projects), "projects.")
		fmt.Println("Info: You can add more projects by running 'tapeless projects add'.")

		fmt.Println()
		fmt.Println("=== Repository Setup ===")

		repos, err := reposService.FetchAndUpdateRepositories(projects)

		if err != nil {
			fmt.Println("Error reading repositories:", err)
			return
		}

		if len(repos) == 0 {
			fmt.Println("No local repositories configured.")
			fmt.Println("Next step: Run 'tapeless repos add' to add a repository.")
			return
		}

		fmt.Println("Found", len(repos), "repositories.")
		fmt.Println("Info: You can add more repositories by running 'tapeless repos add'.")

		fmt.Println()
		fmt.Println("=== Summary ===")
		fmt.Println("You are all set up and ready to go!")
		fmt.Println("Run 'tapeless sync' to push your commits to tapeless.")
		fmt.Println("Generate time entries by running 'tapeless time generate'.")
		fmt.Println("Add manual time entries by running 'tapeless time add'.")
		fmt.Println()
		fmt.Println("We hope you enjoy using Tapeless!")

	},
}
