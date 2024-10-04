package projectsService

import (
	"encoding/json"
	"fmt"
	"io"

	"tapeless.app/tapeless-cli/env"
	"tapeless.app/tapeless-cli/util"
)

func CreateProject(request ProjectsCreateRequest) (Project, error) {
	var Project Project
	err := util.MakeAuthRequestAndParseResponse("POST", env.ApiURL+"/projects", request, &Project)

	return Project, err
}

func DeleteProject(projectId int) error {
	_, err := util.MakeAuthRequest("DELETE", fmt.Sprintf("%s/projects/%d", env.ApiURL, projectId), nil)
	return err
}

func FetchProjects() ([]Project, error) {
	resp, err := util.MakeAuthRequest("GET", env.ApiURL+"/projects", nil)

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

func FilterProjectsById(projectId int, projects *[]Project) (Project, error) {
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
