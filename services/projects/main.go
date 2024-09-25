package projectsService

import (
	"encoding/json"
	"io"

	"github.com/spf13/viper"
	"tapeless.app/tapeless-cli/util"
)

func getProjects() ([]ProjectData, error) {
	resp, err := util.MakeRequest("GET", "http://localhost:4000/cli/projects", nil)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var projectsData []ProjectData

	err = json.Unmarshal(body, &projectsData)

	return projectsData, err

}

func persistProjects(projectsData []ProjectData) (map[int]ProjectData, error) {
	projectsMap := make(map[int]ProjectData)
	for _, project := range projectsData {
		projectsMap[project.Id] = project
	}
	viper.Set("projects", projectsMap)
	err := viper.WriteConfig()
	return projectsMap, err
}

func SyncProjects() (map[int]ProjectData, error) {
	projectsData, err := getProjects()

	if err != nil {
		return nil, err
	}

	projectMap, err := persistProjects(projectsData)

	if err != nil {
		return nil, err
	}

	return projectMap, nil
}
