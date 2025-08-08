package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/models"
)

func getLatestCommit(owner, repo string) (*models.GitCommit, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", owner, repo)
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "majesticcoding.com")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var commits []struct {
		Commit struct {
			Message string `json:"message"`
			Author  struct {
				Date string `json:"date"`
			} `json:"author"`
		} `json:"commit"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &commits); err != nil {
		return nil, err
	}
	if len(commits) == 0 {
		return nil, fmt.Errorf("no commits found")
	}

	commit := &models.GitCommit{
		CommitDate: commits[0].Commit.Author.Date,
		Message:    commits[0].Commit.Message,
	}
	return commit, nil
}

func GitHashHandler(c *gin.Context) {
	owner := "mattmajestic"
	repo := "majesticcoding.com"
	commit, err := getLatestCommit(owner, repo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get git data"})
		return
	}
	c.JSON(http.StatusOK, commit)
}
