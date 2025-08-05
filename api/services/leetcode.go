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

	// 1. Solved Count
	profileURL := fmt.Sprintf("%s/userProfile/%s", baseURL, username)
	profileResp, err := http.Get(profileURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching profile: %w", err)
	}
	defer profileResp.Body.Close()

	profileBody, err := io.ReadAll(profileResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading profile body: %w", err)
	}
	var profile struct {
		SubmitStatsGlobal struct {
			AcSubmissionNum []struct {
				Count int `json:"count"`
			} `json:"acSubmissionNum"`
		} `json:"submitStatsGlobal"`
	}
	if err := json.Unmarshal(profileBody, &profile); err != nil {
		return nil, fmt.Errorf("error decoding profile JSON: %w\nRaw: %s", err, string(profileBody))
	}
	if len(profile.SubmitStatsGlobal.AcSubmissionNum) > 0 {
		summary.SolvedCount = profile.SubmitStatsGlobal.AcSubmissionNum[0].Count
	}

	// 2. Contest Ranking
	rankingURL := fmt.Sprintf("%s/userContestRankingInfo/%s", baseURL, username)
	rankingResp, err := http.Get(rankingURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching ranking: %w", err)
	}
	defer rankingResp.Body.Close()

	rankingBody, err := io.ReadAll(rankingResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading ranking body: %w", err)
	}
	var ranking struct {
		Rating int `json:"rating"`
	}
	if err := json.Unmarshal(rankingBody, &ranking); err != nil {
		return nil, fmt.Errorf("error decoding ranking JSON: %w\nRaw: %s", err, string(rankingBody))
	}
	summary.Ranking = ranking.Rating

	// 3. Language Stats
	langURL := fmt.Sprintf("%s/languageStats?username=%s", baseURL, username)
	langResp, err := http.Get(langURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching languages: %w", err)
	}
	defer langResp.Body.Close()

	langBody, err := io.ReadAll(langResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading lang body: %w", err)
	}
	var langs []struct {
		LanguageName string `json:"languageName"`
	}
	if err := json.Unmarshal(langBody, &langs); err != nil {
		return nil, fmt.Errorf("error decoding language JSON: %w\nRaw: %s", err, string(langBody))
	}
	if len(langs) > 0 {
		summary.MainLanguage = langs[0].LanguageName
	}

	return summary, nil
}
