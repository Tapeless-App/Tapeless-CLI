package prompts

import (
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

	items := make([]RepoSelect, 0)

	for _, repository := range repositories {
		project, err := projectService.FilterProjectsById(repository.ProjectId, &projects)

		if err != nil {
			project.Name = "Unknown"
		}

		items = append(items, RepoSelect{
			Id:          repository.GitConfigId,
			Name:        repository.Name,
			Path:        repository.Path,
			ProjectName: project.Name,
			ProjectId:   repository.ProjectId,
		})
	}

	// sort the items by project name
	sort.Slice(items, func(i, j int) bool {
		return items[i].ProjectId < items[j].ProjectId
	})

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

	items := []projectService.Project{}

	for _, project := range projects {

		items = append(items, projectService.Project{
			Id:           project.Id,
			Name:         project.Name,
			LastSync:     util.DateTimeToDateStr(project.LastSync),
			ProjectStart: util.DateTimeToDateStr(project.ProjectStart),
			ProjectEnd:   util.DateTimeToDateStr(project.ProjectEnd),
		})
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
