package sync

type Commit struct {
	CommitHash string   `json:"commit_hash"`
	Date       string   `json:"date"`
	Message    string   `json:"message"`
	Branches   []string `json:"branches"`
}
