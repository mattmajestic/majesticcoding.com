package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"majesticcoding.com/api/models"
)

const baseURL = "https://alfa-leetcode-api.onrender.com"

func FetchLeetCodeStats(username string) (*models.LeetCodeStats, error) {
	summary := &models.LeetCodeStats{Username: username}
	summary.Languages = "Python | SQL | Go"

	profileURL := fmt.Sprintf("%s/userProfile/%s", baseURL, username)
	resp, err := http.Get(profileURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching profile: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	// Unmarshal only the fields we need
	var parsed struct {
		TotalSolved int `json:"totalSolved"`
		Ranking     int `json:"ranking"`
	}

	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("error decoding profile JSON: %w\nRaw: %s", err, string(body))
	}

	summary.SolvedCount = parsed.TotalSolved
	summary.Ranking = parsed.Ranking

	return summary, nil
}
