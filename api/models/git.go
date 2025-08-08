package models

type GitCommit struct {
	CommitDate string `json:"commit_date"`
	Message    string `json:"message"`
}
