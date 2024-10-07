package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"tapeless.app/tapeless-cli/cmd"
	"tapeless.app/tapeless-cli/env"
	authService "tapeless.app/tapeless-cli/services/auth"
	projectsService "tapeless.app/tapeless-cli/services/projects"
	reposService "tapeless.app/tapeless-cli/services/repos"
	syncService "tapeless.app/tapeless-cli/services/sync"
	versionService "tapeless.app/tapeless-cli/services/version"
	"tapeless.app/tapeless-cli/util"
)

func init() {
	cmd.RootCmd.AddCommand(SyncCmd)
	SyncCmd.Flags().BoolVarP(&includeInactiveFlag, "include-completed", "i", false, "Sync repositories all projects, even if they have been completed for more than 30 days")
}

var (
	includeInactiveFlag bool
	SyncCmd             = &cobra.Command{
		Use:   "sync",
		Short: "Sync the commits from your repositories with Tapeless, will sync all running projects or that have ended within the last 30 days",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			authService.EnsureValidSession()
		},
		Run: func(cmd *cobra.Command, args []string) {

			versionService.CheckLatestVersion()

			projects, err := projectsService.FetchProjects()

			if err != nil {
				fmt.Println("Error reading projects:", err)
				return
			}

			if len(projects) == 0 {
				fmt.Println("No projects found - add a project first using \"tapeless projects add\"")
				return
			}

			activeProjects := make([]projectsService.Project, 0)

			for _, project := range projects {

				if project.ProjectEnd == "" {
					activeProjects = append(activeProjects, project)
					continue
				}

				endDate, err := time.Parse("2006-01-02T15:04:05.000Z", project.ProjectEnd)

				if err != nil {
					fmt.Printf("Error parsing project end date for project %s: %s\n", project.Name, err.Error())
					fmt.Println("Skipping project")
					continue
				}

				if includeInactiveFlag || endDate.After(time.Now().AddDate(0, 0, -30)) {
					activeProjects = append(activeProjects, project)
				}

			}

			// Ensure that remote gitConfigs and local repositories are in sync
			repositories, err := reposService.FetchAndUpdateRepositories(activeProjects)

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

				activeProject, err := projectsService.FilterProjectsById(repo.ProjectId, &activeProjects)

				if err != nil {
					fmt.Println("Error finding project for repository:", repo.Name, err)
					continue
				}

				commits, err := syncService.GetLocalCommitListForRepo(*repo, activeProject)

				if err != nil {
					fmt.Println("Error getting commit list for repository:", repo.Name, err)
					continue
				}

				if len(commits) == 0 {
					fmt.Println("No new commits found for repository:", repo.Name)
					continue
				} else {
					fmt.Println("Found", len(commits), "new commits for repository:", repo.Name)
				}

				// If repo.LastSync is empty, ask the user before proceeding using promptui
				if repo.LatestSync == "" {
					prompt := promptui.Prompt{
						Label:     fmt.Sprintf("Repository %s has never been synced before. Do you want to proceed with syncing these commits?", repo.Name),
						Default:   "Y",
						IsConfirm: true,
					}

					_, err := prompt.Run()

					if err != nil {
						fmt.Println("Skipping repository:", repo.Name)
						continue
					}
				}

				batches := splitCommitsIntoBatches(commits, 250)

				batchCount := len(batches)

				if batchCount > 1 {
					fmt.Println("Commit list split into", batchCount, "batches for repository:", repo.Name)
				}

				for i, batch := range batches {
					// Convert the batch of commits to JSON
					jsonOutput, err := json.Marshal(batch)
					if err != nil {
						fmt.Println("Error marshaling JSON for batch", i+1, ":", err)
						continue
					}

					uploadUrl := fmt.Sprintf("%s/projects/%d/gitConfigs/%d/commits", env.ApiURL, repo.ProjectId, repo.GitConfigId)

					_, err = util.MakeAuthRequest("POST", uploadUrl, bytes.NewBuffer(jsonOutput))

					if err != nil {
						fmt.Println("Error uploading batch", i+1, "of commits:", err)
						continue
					}

					if batchCount > 1 {
						fmt.Println("Batch", i+1, "of commits uploaded successfully for repository:", repo.Name)
					}

				}

				fmt.Println("Commits uploaded successfully for repository:", repo.Name)

				repo.LatestSync = time.Now().Format("2006-01-02T15:04:05")

			}

			fmt.Println("All repositories synced successfully")

			viper.Set("repositories", repositories)
			viper.WriteConfig()

		},
	}
)

// Function to split commits into batches
func splitCommitsIntoBatches(commits []syncService.LocalCommit, batchSize int) [][]syncService.LocalCommit {
	var batches [][]syncService.LocalCommit
	for batchSize < len(commits) {
		batches = append(batches, commits[:batchSize])
		commits = commits[batchSize:]
	}
	if len(commits) > 0 {
		batches = append(batches, commits)
	}
	return batches
}
