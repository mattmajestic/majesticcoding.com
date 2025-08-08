package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/services"
)

// dummy provider handlers map
var statsHandlers = map[string]func(c *gin.Context){
	"youtube":  getYouTubeStats,
	"github":   getGithubStats,
	"twitch":   getTwitchStats,
	"leetcode": getLeetCodeStats,
}

// StatsRouter godoc
// @Summary Get stats from a provider
// @Description Returns stats for the given provider (youtube, github, twitch, leetcode)
// @Tags Stats
// @Param provider path string true "Stats Provider"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /stats/{provider} [get]
func StatsRouter(c *gin.Context) {
	provider := c.Param("provider")

	handler, exists := statsHandlers[provider]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "provider not supported"})
		return
	}

	handler(c)
}

// YouTube handler for API Data
func getYouTubeStats(c *gin.Context) {
	stats, err := services.FetchYouTubeStats()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GitHub Handler for API Data
func getGithubStats(c *gin.Context) {
	username := "mattmajestic"

	stats, err := services.FetchGitHubStats(username)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// Twitch Handler for API Data
func getTwitchStats(c *gin.Context) {
	username := "MajesticCodingTwitch"

	stats, err := services.FetchTwitchStats(username)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// Dummy Leetcode handler
// func getLeetCodeStats(c *gin.Context) {
// 	dummy := models.LeetCodeStats{
// 		Username:     "mattmajestic",
// 		MainLanguage: "Python & Go",
// 		SolvedCount:  355,
// 		Ranking:      285162,
// 	}
// 	c.JSON(http.StatusOK, dummy)
// }

func getLeetCodeStats(c *gin.Context) {
	username := "mattmajestic"

	stats, err := services.FetchLeetCodeStats(username)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
