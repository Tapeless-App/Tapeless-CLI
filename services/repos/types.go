package reposService

type LocalRepositoryConfig struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	AuthorEmail string `json:"authorEmail"`
	OriginUrl   string `json:"originUrl"`
}

type Repository struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	LatestSync  string `json:"latestSync"`
	ProjectId   int    `json:"projectId"`
	GitConfigId int    `json:"gitConfigId"`
	AuthorEmail string `json:"authorEmail"`
	OriginUrl   string `json:"originUrl"`
}

type GitConfigRequest struct {
	Name          string `json:"name"`
	AuthorEmail   string `json:"authorEmail"`
	RepositoryUrl string `json:"repositoryUrl"`
	ProjectId     int    `json:"projectId"`
}

type GitConfigResponse struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	AuthorEmail   string `json:"authorEmail"`
	RepositoryUrl string `json:"repositoryUrl"`
	ProjectId     int    `json:"projectId"`
}
