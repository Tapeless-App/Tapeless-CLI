package projectsService

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/viper"
	"tapeless.app/tapeless-cli/util"
)

func CreateProject(request ProjectsCreateRequest) (Project, error) {
	var Project Project
	err := util.MakeRequestAndParseResponse("POST", "http://localhost:4000/cli/projects", request, &Project)

	return Project, err
}

func getProjects() ([]Project, error) {
	resp, err := util.MakeRequest("GET", "http://localhost:4000/cli/projects", nil)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var projects []Project

	err = json.Unmarshal(body, &projects)

	return projects, err
}

func DeleteProject(projectId int) error {
	_, err := util.MakeRequest("DELETE", fmt.Sprintf("http://localhost:4000/cli/projects/%d", projectId), nil)
	return err
}

func persistProjects(projects []Project) ([]Project, error) {
	viper.Set("projects", projects)
	err := viper.WriteConfig()
	return projects, err
}

func SyncProjects() ([]Project, error) {
	projects, err := getProjects()

	if err != nil {
		return nil, err
	}

	_, err = persistProjects(projects)

	if err != nil {
		return nil, err
	}

	return projects, nil
}

func FindProjectById(projectId int, projects *[]Project) (Project, error) {
	if projects == nil {
		return Project{}, fmt.Errorf("projects is nil")
	}

	for _, project := range *projects {
		if project.Id == projectId {
			return project, nil
		}
	}

	return Project{}, fmt.Errorf("project not found")
}
