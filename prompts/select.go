package prompts

import (
	"fmt"
	"slices"
	"sort"

	"github.com/manifoldco/promptui"
	projectService "tapeless.app/tapeless-cli/services/projects"
	reposService "tapeless.app/tapeless-cli/services/repos"

	timeService "tapeless.app/tapeless-cli/services/time"
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
func GetProjectIdPrompt(label string, projectIdFlag int, projects []projectService.Project) (projectService.Project, error) {
	return GetProjectIdPromptWithDefault(label, projectIdFlag, projects, -1)
}

/**
 * Get the project ID for the repository
 * Will use the flag if it is set, otherwise prompt the user with a list of projects
 */
func GetProjectIdPromptWithDefault(label string, projectIdFlag int, projects []projectService.Project, defaultProjectId int) (projectService.Project, error) {

	if projectIdFlag != -1 {
		return projectService.FilterProjectsById(projectIdFlag, &projects)
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

	defaultIndex := slices.IndexFunc(items, func(p projectService.Project) bool {
		return p.Id == defaultProjectId
	})

	if defaultIndex == -1 {
		defaultIndex = 0
	}

	prompt := promptui.Select{
		Templates: templates,
		Label:     label,
		Items:     items,
		Size:      len(items),
		CursorPos: defaultIndex,
	}

	index, _, err := prompt.Run()

	if err != nil {
		return projectService.Project{}, err
	}

	return projectService.FilterProjectsById(items[index].Id, &projects)
}

func SelectTimeEntryPrompt(label string, timeEntries []timeService.TimeEntry) (timeService.TimeEntry, error) {

	type TimeEntrySelect struct {
		timeService.TimeEntry
		ShortDescription string
	}

	items := make([]TimeEntrySelect, 0)

	templates := &promptui.SelectTemplates{
		Label:    `{{ . }}:`,
		Active:   "> {{ .ShortDescription | cyan }} | {{.Hours}} hours",
		Inactive: "  {{ .ShortDescription | cyan }} | {{.Hours}} hours",
		Selected: "{{ .ShortDescription }}",
		Details: `
--------- Time Entry: {{ .ShortDescription }} ----------
{{ "Id:" | faint }}	{{ .Id }}
{{ "Date:" | faint }}	{{ .Date }}
{{ "Hours:" | faint }}	{{ .Hours }}
{{ "Description:" | faint }}	{{ .Description }}

 `}

	for _, entry := range timeEntries {

		shortDescription := entry.Description

		if len(entry.Description) > 20 {
			shortDescription = entry.Description[:20] + "..."
		}

		items = append(items, TimeEntrySelect{
			TimeEntry:        entry,
			ShortDescription: shortDescription,
		})
	}

	prompt := promptui.Select{
		Templates: templates,
		Label:     label,
		Items:     items,
		Size:      len(items),
	}

	index, _, err := prompt.Run()

	if err != nil {
		return timeService.TimeEntry{}, err
	}

	return items[index].TimeEntry, nil
}
