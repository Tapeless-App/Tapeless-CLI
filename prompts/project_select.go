package prompts

import (
	"github.com/manifoldco/promptui"
	projectService "tapeless.app/tapeless-cli/services/projects"
	"tapeless.app/tapeless-cli/util"
)

/**
 * Get the project ID for the repository
 * Will use the flag if it is set, otherwise prompt the user with a list of projects
 */
func GetProjectIdPrompt(label string, projectIdFlag int, projects map[int]projectService.ProjectData) (int, error) {

	if projectIdFlag != -1 {
		return projectIdFlag, nil
	}

	items := []projectService.ProjectData{}

	for _, project := range projects {

		items = append(items, projectService.ProjectData{
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
