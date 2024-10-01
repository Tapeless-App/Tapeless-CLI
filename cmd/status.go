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
	Use:   "status",
	Short: "Get your Tapeless setup and configuration status",
	Run: func(cmd *cobra.Command, args []string) {

		token := viper.GetString("token")

		fmt.Println("=== Login Status ===")

		if token == "" {
			fmt.Println("Not logged in.")
			fmt.Println("Next step: Run `tapeless login` login to tapeless.")
			return
		}

		isExpired, err := util.IsJWTExpired(token)

		if err != nil {
			fmt.Println("Error verifying access token:", err)
			fmt.Println("Next step: Try to run `tapeless login` - if that doesn't work try deleting the config file and starting over.")
			return
		}

		if isExpired {
			fmt.Println("JWT token expired.")
			fmt.Println("Next step: Run `tapeless login` login to tapeless.")
			return
		}

		fmt.Println("Logged in.")

		fmt.Println("=== Project Setup ===")

		projects, err := projectsService.SyncProjects()

		if err != nil {
			fmt.Println("Error reading projects:", err)
			return
		}

		if len(projects) == 0 {
			fmt.Println("No projects found.")
			fmt.Println("Next step: Run `tapeless projects add` to add a project.")
			return
		}

		fmt.Println("Found", len(projects), "projects.")

		fmt.Println("=== Repository Setup ===")

		repos, err := reposService.SyncRepositories(projects)

		if err != nil {
			fmt.Println("Error reading repositories:", err)
			return
		}

		if len(repos) == 0 {
			fmt.Println("No local repositories configured.")
			fmt.Println("Next step: Run `tapeless repos add` to add a repository.")
			return
		}

		fmt.Println("Found", len(repos), "repositories.")

		fmt.Println("Next step: Run `tapeless sync` to sync your repositories with Tapeless.")

	},
}
