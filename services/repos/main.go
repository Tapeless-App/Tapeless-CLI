package reposService

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"tapeless.app/tapeless-cli/util"
)

func GetPersistedRepositories() ([]Repository, error) {
	currentRepositories := make([]Repository, 0)

	err := viper.UnmarshalKey("repositories", &currentRepositories)

	return currentRepositories, err
}

func CreateGitConfig(projectId int, localRepo LocalRepositoryConfig) (GitConfigResponse, error) {

	gitConfigRequest := GitConfigRequest{
		Name:          localRepo.Name,
		AuthorEmail:   localRepo.AuthorEmail,
		RepositoryUrl: localRepo.OriginUrl,
		ProjectId:     projectId,
	}

	gitConfigResponse := GitConfigResponse{}

	err := util.MakeRequestAndParseResponse("POST", fmt.Sprintf("http://localhost:4000/cli/projects/%d/gitConfigs", projectId), gitConfigRequest, &gitConfigResponse)

	return gitConfigResponse, err
}

func DeleteGitConfig(projectId int, gitConfigId int) error {
	_, err := util.MakeRequest("DELETE", fmt.Sprintf("http://localhost:4000/cli/projects/%d/gitConfigs/%d", projectId, gitConfigId), nil)
	return err
}

func GetLocalRepositoryConfig(path string) (LocalRepositoryConfig, error) {

	if path == "" {
		var error error

		path, error = os.Getwd()

		if error != nil {
			return LocalRepositoryConfig{}, error
		}
	}

	path, err := filepath.Abs(path)

	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return LocalRepositoryConfig{}, err
	}

	inGitRepoCmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	inGitRepoOutput, err := inGitRepoCmd.Output()

	repoNameSegments := strings.Split(path, string(filepath.Separator))

	repoName := repoNameSegments[len(repoNameSegments)-1]

	if err != nil {
		if strings.Contains(err.Error(), "exit status 128") {
			fmt.Println("Not in git LocalRepositoryConfig!")
			return LocalRepositoryConfig{}, err
		}
		fmt.Println("Error Str:", err)
		return LocalRepositoryConfig{}, err
	}

	// Technically unnecessary, will return error if not in repo, but just to be safe
	if !strings.Contains(string(inGitRepoOutput), "true") {
		fmt.Println("Not in git LocalRepositoryConfig!")
		return LocalRepositoryConfig{}, err
	}

	gitAuthorCmd := exec.Command("git", "config", "--get", "user.email")

	gitAuthorOutput, err := gitAuthorCmd.Output()

	if err != nil {
		fmt.Println("Error fetching git author:", err)
		return LocalRepositoryConfig{}, err
	}

	gitOriginCmd := exec.Command("git", "config", "--get", "remote.origin.url")

	gitOriginOutput, err := gitOriginCmd.Output()

	if err != nil {
		fmt.Println("Error fetching git origin:", err)
		return LocalRepositoryConfig{}, err
	}

	return LocalRepositoryConfig{
		Name:        repoName,
		Path:        path,
		AuthorEmail: strings.TrimSpace(string(gitAuthorOutput)),
		OriginUrl:   strings.TrimSpace(string(gitOriginOutput)),
	}, nil

}

// func GetCommitsJsonStringForRepository(repo Repository) (string, error) {

// }
