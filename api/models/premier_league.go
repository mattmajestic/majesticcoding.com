package models

import "time"

type PLMatch struct {
	ID        int       `json:"id"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"`
	Matchday  int       `json:"matchday"`
	HomeTeam  PLTeam    `json:"home_team"`
	AwayTeam  PLTeam    `json:"away_team"`
	Score     PLScore   `json:"score"`
}

type PLTeam struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Crest string `json:"crest"`
}

type PLScore struct {
	Winner   string `json:"winner"`
	Duration string `json:"duration"`
	FullTime PLResult `json:"full_time"`
	HalfTime PLResult `json:"half_time"`
}

type PLResult struct {
	Home *int `json:"home"`
	Away *int `json:"away"`
}

type PLScheduleResponse struct {
	Matches []PLMatch `json:"matches"`
}