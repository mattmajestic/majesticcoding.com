package models

type LeetCodeStats struct {
	Username     string `json:"username"`
	SolvedCount  int    `json:"solved_count"`
	Ranking      int    `json:"ranking"`
	MainLanguage string `json:"main_language"`
}
