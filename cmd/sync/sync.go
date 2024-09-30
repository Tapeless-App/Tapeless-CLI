package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"tapeless.app/tapeless-cli/cmd"
	"tapeless.app/tapeless-cli/env"
	projectsService "tapeless.app/tapeless-cli/services/projects"
	reposService "tapeless.app/tapeless-cli/services/repos"
	syncService "tapeless.app/tapeless-cli/services/sync"
	"tapeless.app/tapeless-cli/util"
)

func init() {
	cmd.RootCmd.AddCommand(syncCmd)
}

var (
	syncCmd = &cobra.Command{
		Use:   "sync",
		Short: "Sync the commits from your repositories with Tapeless",
		Run: func(cmd *cobra.Command, args []string) {

			projects, err := projectsService.SyncProjects()

			if err != nil {
				fmt.Println("Error reading projects:", err)
				return
			}

			// Ensure that remote gitConfigs and local repositories are in sync
			repositories, err := reposService.SyncRepositories(projects)

			if err != nil {
				fmt.Println("Error reading repositories:", err)
				return
			}

			if len(repositories) == 0 {
				fmt.Println("No repositories found - add a repository first using \"tapeless repos add\"")
				return
			}

			for repoIndex := range repositories {

				repo := &repositories[repoIndex]
				project := projects[repo.ProjectId]

				commits, err := syncService.GetCommitListForRepo(*repo, project)

				if err != nil {
					fmt.Println("Error getting commit list for repository:", repo.Name, err)
					return
				}

				if len(commits) == 0 {
					fmt.Println("No new commits found for repository:", repo.Name)
					continue
				} else {
					fmt.Println("Found", len(commits), "new commits for repository:", repo.Name)
				}

				// Convert the list of commits to JSON
				jsonOutput, err := json.Marshal(commits)
				if err != nil {
					fmt.Println("Error marshaling JSON:", err)
					return
				}

				uploadUrl := fmt.Sprintf("%s/projects/%d/gitConfigs/%d/commits", env.ApiURL, repo.ProjectId, repo.GitConfigId)

				_, err = util.MakeRequest("POST", uploadUrl, bytes.NewBuffer(jsonOutput))

				if err != nil {
					fmt.Println("Error uploading commits:", err)
					return
				}

				fmt.Println("Commits uploaded successfully for repo,", repo.Name, "- updating latest sync time")

				repo.LatestSync = time.Now().Format("2006-01-02T15:04:05")

			}

			viper.Set("repositories", repositories)
			viper.WriteConfig()

		},
	}
)
