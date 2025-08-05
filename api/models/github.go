package models

type GitHubStats struct {
	Username      string `json:"username"`
	PublicRepos   int    `json:"public_repos"`
	Followers     int    `json:"followers"`
	StarsReceived int    `json:"stars_received"`
}
