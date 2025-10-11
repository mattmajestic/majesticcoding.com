package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"majesticcoding.com/api/models"
)

const baseURL = "https://alfa-leetcode-api.onrender.com"

var ErrRateLimited = errors.New("leetcode API rate limited")

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

	// Check for rate limiting message
	bodyStr := string(body)
	if strings.Contains(bodyStr, "Too many request") {
		return nil, ErrRateLimited
	}

	// Unmarshal only the fields we need
	var parsed struct {
		TotalSolved int `json:"totalSolved"`
		Ranking     int `json:"ranking"`
	}

	if err := json.Unmarshal(body, &parsed); err != nil {
		// Check if the error message also indicates rate limiting
		if strings.Contains(bodyStr, "Too many request") {
			return nil, ErrRateLimited
		}
		return nil, fmt.Errorf("error decoding profile JSON: %w\nRaw: %s", err, bodyStr)
	}

	summary.SolvedCount = parsed.TotalSolved
	summary.Ranking = parsed.Ranking

	return summary, nil
}
