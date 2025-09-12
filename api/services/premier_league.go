package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"majesticcoding.com/api/models"
)

func FetchPLSchedule() ([]models.PLMatch, error) {
	apiKey := os.Getenv("EPL_TOKEN")
	if apiKey == "" {
		return nil, fmt.Errorf("EPL_TOKEN not found")
	}

	// Get matches for the current week (next 7 days)
	today := time.Now()
	from := today.Format("2006-01-02")
	to := today.AddDate(0, 0, 7).Format("2006-01-02")

	url := fmt.Sprintf("https://api.football-data.org/v4/competitions/PL/matches?dateFrom=%s&dateTo=%s", from, to)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("request creation failed: %w", err)
	}

	req.Header.Set("X-Auth-Token", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var apiResponse struct {
		Matches []struct {
			ID       int    `json:"id"`
			UTCDate  string `json:"utcDate"`
			Status   string `json:"status"`
			Matchday int    `json:"matchday"`
			HomeTeam struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Crest string `json:"crest"`
			} `json:"homeTeam"`
			AwayTeam struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Crest string `json:"crest"`
			} `json:"awayTeam"`
			Score struct {
				Winner   *string `json:"winner"`
				Duration string  `json:"duration"`
				FullTime struct {
					Home *int `json:"home"`
					Away *int `json:"away"`
				} `json:"fullTime"`
				HalfTime struct {
					Home *int `json:"home"`
					Away *int `json:"away"`
				} `json:"halfTime"`
			} `json:"score"`
		} `json:"matches"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	matches := make([]models.PLMatch, len(apiResponse.Matches))
	for i, match := range apiResponse.Matches {
		date, _ := time.Parse(time.RFC3339, match.UTCDate)
		
		winner := ""
		if match.Score.Winner != nil {
			winner = *match.Score.Winner
		}

		matches[i] = models.PLMatch{
			ID:       match.ID,
			Date:     date,
			Status:   match.Status,
			Matchday: match.Matchday,
			HomeTeam: models.PLTeam{
				ID:   match.HomeTeam.ID,
				Name: match.HomeTeam.Name,
				Crest: match.HomeTeam.Crest,
			},
			AwayTeam: models.PLTeam{
				ID:   match.AwayTeam.ID,
				Name: match.AwayTeam.Name,
				Crest: match.AwayTeam.Crest,
			},
			Score: models.PLScore{
				Winner:   winner,
				Duration: match.Score.Duration,
				FullTime: models.PLResult{
					Home: match.Score.FullTime.Home,
					Away: match.Score.FullTime.Away,
				},
				HalfTime: models.PLResult{
					Home: match.Score.HalfTime.Home,
					Away: match.Score.HalfTime.Away,
				},
			},
		}
	}

	return matches, nil
}