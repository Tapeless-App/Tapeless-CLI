package syncService

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"tapeless.app/tapeless-cli/env"
	projectsService "tapeless.app/tapeless-cli/services/projects"
	reposService "tapeless.app/tapeless-cli/services/repos"
	"tapeless.app/tapeless-cli/util"
)

func GetLocalCommitListForRepo(repo reposService.Repository, project projectsService.Project) ([]LocalCommit, error) {

	authorFlag := fmt.Sprintf("--author=%s", repo.AuthorEmail)
	sinceFlag := fmt.Sprintf("--since=%s", repo.LatestSync)

	var commitList []LocalCommit

	if repo.LatestSync == "" {
		sinceFlag = fmt.Sprintf("--since=%s", util.DateTimeToDateStr(project.ProjectStart))
	}

	gitCommitsCmd := exec.Command("git", "log", "--all", authorFlag, sinceFlag, "--date=format:%Y-%m-%d %H:%M:%S", "--pretty=format:%H,%ad,%s")

	gitCommitsCmd.Dir = repo.Path
	var gitCommitsCmdOut bytes.Buffer
	gitCommitsCmd.Stdout = &gitCommitsCmdOut

	err := gitCommitsCmd.Run()

	if err != nil {
		return nil, err
	}

	// Split the output by line
	commits := strings.Split(gitCommitsCmdOut.String(), "\n")

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
		branchesCmd.Dir = repo.Path
		var branchesOut bytes.Buffer
		branchesCmd.Stdout = &branchesOut
		err = branchesCmd.Run()
		if err != nil {
			return nil, err
		}

		// Process the branches output, clean and split
		branches := strings.Fields(branchesOut.String())
		for i, branch := range branches {
			branches[i] = strings.TrimSpace(branch)
		}

		// Create a commit entry
		commit := LocalCommit{
			CommitHash: commitHash,
			Date:       date,
			Message:    message,
			Branches:   branches,
		}
		commitList = append(commitList, commit)
	}

	return commitList, nil
}

func FetchCommitsForProjectAndDate(projectId int, date string) ([]Commit, error) {
	var commits []Commit
	err := util.MakeAuthRequestAndParseResponse("GET", fmt.Sprintf("%s/projects/%d/commits?date=%s", env.ApiURL, projectId, date), nil, &commits)
	return commits, err
}
