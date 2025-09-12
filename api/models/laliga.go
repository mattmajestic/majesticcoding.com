package models

import "time"

type LaLigaMatch struct {
	ID        int       `json:"id"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"`
	Matchday  int       `json:"matchday"`
	HomeTeam  LaLigaTeam `json:"home_team"`
	AwayTeam  LaLigaTeam `json:"away_team"`
	Score     LaLigaScore `json:"score"`
}

type LaLigaTeam struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Crest string `json:"crest"`
}

type LaLigaScore struct {
	Winner   string `json:"winner"`
	Duration string `json:"duration"`
	FullTime LaLigaResult `json:"full_time"`
	HalfTime LaLigaResult `json:"half_time"`
}

type LaLigaResult struct {
	Home *int `json:"home"`
	Away *int `json:"away"`
}

type LaLigaScheduleResponse struct {
	Matches []LaLigaMatch `json:"matches"`
}