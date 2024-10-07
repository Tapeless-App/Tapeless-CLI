package prompts

import (
	"fmt"
	"slices"
	"sort"

	"github.com/manifoldco/promptui"
	projectService "tapeless.app/tapeless-cli/services/projects"
	reposService "tapeless.app/tapeless-cli/services/repos"
	"tapeless.app/tapeless-cli/util"
)

type RepoSelect struct {
	Id          int
	Name        string
	Path        string
	ProjectName string
	ProjectId   int
}

func GetRepositoryPrompt(label string, repositories []reposService.Repository, projects []projectService.Project) (*reposService.Repository, error) {

	activeItems := make([]RepoSelect, 0)
	completedItems := make([]RepoSelect, 0)

	for _, repository := range repositories {
		project, err := projectService.FilterProjectsById(repository.ProjectId, &projects)

		if err != nil {
			project.Name = "Unknown"
		}

		isCompleted, err := util.IsDateInPast("2006-01-02T15:04:05.000Z", project.ProjectEnd)

		if err != nil || !isCompleted {
			activeItems = append(activeItems, RepoSelect{
				Id:          repository.GitConfigId,
				Name:        repository.Name,
				Path:        repository.Path,
				ProjectName: project.Name,
				ProjectId:   repository.ProjectId,
			})
		} else {
			completedItems = append(completedItems, RepoSelect{
				Id:          repository.GitConfigId,
				Name:        repository.Name,
				Path:        repository.Path,
				ProjectName: fmt.Sprintf("%s [COMPLETED]", project.Name),
				ProjectId:   repository.ProjectId,
			})
		}

	}

	sort.Slice(activeItems, func(i, j int) bool {
		return activeItems[i].ProjectId < activeItems[j].ProjectId
	})

	sort.Slice(completedItems, func(i, j int) bool {
		return completedItems[i].ProjectId < completedItems[j].ProjectId
	})

	items := slices.Concat(activeItems, completedItems)

	templates := &promptui.SelectTemplates{
		Label:    `{{ . }}:`,
		Active:   "> {{ .Name | cyan }}: {{ .Path }} (Linked project: {{ .ProjectName }})",
		Inactive: "  {{ .Name | cyan }}: {{ .Path }} (Linked project: {{ .ProjectName }})",
		Selected: "{{ .Name }}",
	}

	prompt := promptui.Select{
		Templates: templates,
		Label:     label,
		Items:     items,
		Size:      len(items),
	}

	index, _, err := prompt.Run()

	if err != nil {
		return nil, err
	}

	selectedRepo := &reposService.Repository{}

	for _, repo := range repositories {
		if repo.GitConfigId == items[index].Id {
			selectedRepo = &repo
			break
		}
	}

	return selectedRepo, nil
}

/**
 * Get the project ID for the repository
 * Will use the flag if it is set, otherwise prompt the user with a list of projects
 */
func GetProjectIdPrompt(label string, projectIdFlag int, projects []projectService.Project) (int, error) {

	if projectIdFlag != -1 {
		return projectIdFlag, nil
	}

	activeItems := make([]projectService.Project, 0)

	completedItems := make([]projectService.Project, 0)

	for _, project := range projects {

		isCompleted, err := util.IsDateInPast("2006-01-02T15:04:05.000Z", project.ProjectEnd)

		if err != nil || !isCompleted {
			activeItems = append(activeItems, projectService.Project{
				Id:           project.Id,
				Name:         project.Name,
				LastSync:     util.DateTimeToDateStr(project.LastSync),
				ProjectStart: util.DateTimeToDateStr(project.ProjectStart),
				ProjectEnd:   util.DateTimeToDateStr(project.ProjectEnd),
			})
		} else {
			completedItems = append(completedItems, projectService.Project{
				Id:           project.Id,
				Name:         fmt.Sprintf("%s [COMPLETED]", project.Name),
				LastSync:     util.DateTimeToDateStr(project.LastSync),
				ProjectStart: util.DateTimeToDateStr(project.ProjectStart),
				ProjectEnd:   util.DateTimeToDateStr(project.ProjectEnd),
			})
		}

	}

	templates := &promptui.SelectTemplates{
		Label:    `{{ . }}:`,
		Active:   "> {{ .Name | cyan }} (id: {{ .Id }})",
		Inactive: "  {{ .Name | cyan }} (id: {{ .Id }})",
		Selected: "{{ .Name }}",
		Details: `
--------- Project: {{ .Name }} ----------

{{ "Id:" | faint }}	{{ .Id }}
{{ "Project Start:" | faint }}	{{ .ProjectStart }}
{{ "Project End:" | faint }}	{{ .ProjectEnd }}
{{ "Last Sync:" | faint }}	{{ .LastSync }}
 
`,
	}

	items := slices.Concat(activeItems, completedItems)

	prompt := promptui.Select{
		Templates: templates,
		Label:     label,
		Items:     items,
		Size:      len(items),
	}

	index, _, err := prompt.Run()

	if err != nil {
		return -1, err
	}

	return items[index].Id, nil
}
