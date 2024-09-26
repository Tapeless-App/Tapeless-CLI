package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"tapeless.app/tapeless-cli/cmd"
	projectsService "tapeless.app/tapeless-cli/services/projects"
	reposService "tapeless.app/tapeless-cli/services/repos"
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

			// Sync repositories
			repositories, err := reposService.GetPersistedRepositories()

			if err != nil {
				fmt.Println("Error reading repositories:", err)
				return
			}

			projects, err := projectsService.SyncProjects()

			if err != nil {
				fmt.Println("Error reading projects:", err)
				return
			}

			if len(repositories) == 0 {
				fmt.Println("No repositories found - add a repository first using \"tapeless repos add\"")
				return
			}

			for repoIndex := range repositories {

				repo := &repositories[repoIndex]

				fmt.Println("Syncing repository:", repo.Name)

				authorFlag := fmt.Sprintf("--author=%s", repo.AuthorEmail)
				sinceFlag := fmt.Sprintf("--since=%s", repo.LatestSync)

				if repo.LatestSync == "" {
					sinceFlag = fmt.Sprintf("--since=%s", util.DateTimeToDateStr(projects[repo.ProjectId].ProjectStart))
				}

				gitCommitsCmd := exec.Command("git", "log", "--all", authorFlag, sinceFlag, "--date=format:%Y-%m-%d %H:%M:%S", "--pretty=format:%H,%ad,%s")

				gitCommitsCmd.Dir = repo.Path
				var gitCommitsCmdOut bytes.Buffer
				gitCommitsCmd.Stdout = &gitCommitsCmdOut

				err = gitCommitsCmd.Run()

				if err != nil {
					fmt.Println("Error running git log command:", err)
					return
				}

				// Split the output by line
				commits := strings.Split(gitCommitsCmdOut.String(), "\n")
				var commitList []Commit

				for _, commitLine := range commits {
					// Split each line into commit hash, date, and message
					parts := strings.SplitN(commitLine, ",", 3)
					if len(parts) < 3 {
						continue
					}
					commitHash := parts[0]
					date := parts[1]
					message := parts[2]

					// Find the branches that contain this commit
					branchesCmd := exec.Command("git", "branch", "--contains", commitHash)
					var branchesOut bytes.Buffer
					branchesCmd.Stdout = &branchesOut
					err = branchesCmd.Run()
					if err != nil {
						fmt.Println("Error running git branch:", err)
						return
					}

					// Process the branches output, clean and split
					branches := strings.Fields(branchesOut.String())
					for i, branch := range branches {
						branches[i] = strings.TrimSpace(branch)
					}

					// Create a commit entry
					commit := Commit{
						CommitHash: commitHash,
						Date:       date,
						Message:    message,
						Branches:   branches,
					}
					commitList = append(commitList, commit)
				}

				if len(commitList) == 0 {
					fmt.Println("No new commits found for repository:", repo.Name)
					continue
				} else {
					fmt.Println("Found", len(commitList), "new commits for repository:", repo.Name)
				}

				// Convert the list of commits to JSON
				jsonOutput, err := json.Marshal(commitList)
				if err != nil {
					fmt.Println("Error marshaling JSON:", err)
					return
				}

				uploadUrl := fmt.Sprintf("http://localhost:4000/cli/projects/%d/gitConfigs/%d/commits", repo.ProjectId, repo.GitConfigId)

				_, err = util.MakeRequest("POST", uploadUrl, bytes.NewBuffer(jsonOutput))

				if err != nil {
					fmt.Println("Error uploading commits:", err)
					return
				}

				fmt.Println("Commits uploaded successfully - updating latest sync time")

				repo.LatestSync = time.Now().Format("2006-01-02T15:04:05")

			}

			viper.Set("repositories", repositories)
			viper.WriteConfig()

		},
	}
)
