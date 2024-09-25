package projectsService

type ProjectData struct {
	// The fields should match the JSON structure
	// You can use `json` tags to map the fields correctly
	Id           int    `json:"id"`
	Name         string `json:"name"`
	LastSync     string `json:"lastSync"`
	ProjectStart string `json:"projectStart"`
	ProjectEnd   string `json:"projectEnd"`
}
