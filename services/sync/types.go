package syncService

type LocalCommit struct {
	CommitHash string   `json:"commit_hash"`
	Date       string   `json:"date"`
	Message    string   `json:"message"`
	Branches   []string `json:"branches"`
}

type Commit struct {
	ExternalReferenceId string `json:"externalReferenceId"`
	Date                string `json:"date"`
	AuthorEmail         string `json:"authorEmail"`
	Message             string `json:"message"`
	BranchName          string `json:"branchName"`
}
