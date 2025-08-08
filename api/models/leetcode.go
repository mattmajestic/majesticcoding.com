package models

type LeetCodeStats struct {
	Username    string `json:"username"`
	Languages   string `json:"mainLanguages"`
	SolvedCount int    `json:"totalSolved"`
	Ranking     int    `json:"ranking"`
}
