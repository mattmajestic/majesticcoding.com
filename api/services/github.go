package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"majesticcoding.com/api/models"
)

func FetchGitHubStats(username string) (models.GitHubStats, error) {
	client := &http.Client{}

	// Fetch basic user info
	userURL := fmt.Sprintf("https://api.github.com/users/%s", username)
	reqUser, _ := http.NewRequest("GET", userURL, nil)
	reqUser.Header.Set("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))

	respUser, err := client.Do(reqUser)
	if err != nil {
		return models.GitHubStats{}, err
	}
	defer respUser.Body.Close()

	var userData struct {
		PublicRepos int `json:"public_repos"`
		Followers   int `json:"followers"`
	}
	if err := json.NewDecoder(respUser.Body).Decode(&userData); err != nil {
		return models.GitHubStats{}, err
	}

	// Fetch repos to calculate stars
	reposURL := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100", username)
	reqRepos, _ := http.NewRequest("GET", reposURL, nil)
	reqRepos.Header.Set("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))

	respRepos, err := client.Do(reqRepos)
	if err != nil {
		return models.GitHubStats{}, err
	}
	defer respRepos.Body.Close()

	var repos []struct {
		StargazersCount int `json:"stargazers_count"`
	}
	if err := json.NewDecoder(respRepos.Body).Decode(&repos); err != nil {
		return models.GitHubStats{}, err
	}

	totalStars := 0
	for _, repo := range repos {
		totalStars += repo.StargazersCount
	}

	return models.GitHubStats{
		Username:      username,
		PublicRepos:   userData.PublicRepos,
		Followers:     userData.Followers,
		StarsReceived: totalStars,
	}, nil
}
