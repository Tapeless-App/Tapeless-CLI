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
	authService "tapeless.app/tapeless-cli/services/auth"
	projectsService "tapeless.app/tapeless-cli/services/projects"
	reposService "tapeless.app/tapeless-cli/services/repos"
	syncService "tapeless.app/tapeless-cli/services/sync"
	versionService "tapeless.app/tapeless-cli/services/version"
	"tapeless.app/tapeless-cli/util"
)

func init() {
	cmd.RootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVarP(&includeInactiveFlag, "include-completed", "i", false, "Sync repositories all projects, even if they have been completed for more than 30 days")
}

var (
	includeInactiveFlag bool
	syncCmd             = &cobra.Command{
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
					return
				}

				commits, err := syncService.GetCommitListForRepo(*repo, activeProject)

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

				_, err = util.MakeAuthRequest("POST", uploadUrl, bytes.NewBuffer(jsonOutput))

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
