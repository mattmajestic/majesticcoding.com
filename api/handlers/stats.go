package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/services"
	"majesticcoding.com/db"
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
	cacheKey := "youtube:stats:channel"

	// Try to get from Redis cache first (30 minutes TTL)
	cachedJSON, err := services.RedisGetRawJSON(cacheKey)
	if err == nil && cachedJSON != "" {
		log.Printf("✅ YouTube stats cache HIT")
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, cachedJSON)
		return
	}
	log.Printf("🔍 YouTube stats cache MISS, fetching from API")

	stats, err := services.FetchYouTubeStats()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Cache the stats for 30 minutes (1800 seconds)
	if err := services.RedisSetJSON(cacheKey, stats, 1800); err != nil {
		log.Printf("⚠️ Failed to cache YouTube stats: %v", err)
	} else {
		log.Printf("💾 Cached YouTube stats")
	}

	// Store stats in database
	database := db.GetDB()
	if database != nil {
		channelTitle := ""
		subscribers := 0
		videos := 0
		views := int64(0)

		if val, ok := stats["channelTitle"].(string); ok {
			channelTitle = val
		}
		if val, ok := stats["subscribers"].(string); ok {
			fmt.Sscanf(val, "%d", &subscribers)
		}
		if val, ok := stats["videos"].(string); ok {
			fmt.Sscanf(val, "%d", &videos)
		}
		if val, ok := stats["views"].(string); ok {
			fmt.Sscanf(val, "%d", &views)
		}

		if err := db.InsertYouTubeStats(database, channelTitle, subscribers, videos, views); err != nil {
			log.Printf("❌ Failed to save YouTube stats: %v", err)
		}
	}

	c.JSON(http.StatusOK, stats)
}

// GitHub Handler for API Data
func getGithubStats(c *gin.Context) {
	username := "mattmajestic"
	cacheKey := fmt.Sprintf("github:stats:%s", username)

	// Try to get from Redis cache first (30 minutes TTL)
	cachedJSON, err := services.RedisGetRawJSON(cacheKey)
	if err == nil && cachedJSON != "" {
		log.Printf("✅ GitHub stats cache HIT")
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, cachedJSON)
		return
	}
	log.Printf("🔍 GitHub stats cache MISS, fetching from API")

	stats, err := services.FetchGitHubStats(username)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Cache the stats for 30 minutes (1800 seconds)
	if err := services.RedisSetJSON(cacheKey, stats, 1800); err != nil {
		log.Printf("⚠️ Failed to cache GitHub stats: %v", err)
	} else {
		log.Printf("💾 Cached GitHub stats")
	}

	// Store stats in database
	database := db.GetDB()
	if database != nil {
		if err := db.InsertGitHubStats(database, username, stats.PublicRepos, stats.Followers, 0, stats.StarsReceived); err != nil {
			log.Printf("❌ Failed to save GitHub stats: %v", err)
		}
	}

	c.JSON(http.StatusOK, stats)
}

// Twitch Handler for API Data
func getTwitchStats(c *gin.Context) {
	username := "MajesticCodingTwitch"
	cacheKey := fmt.Sprintf("twitch:stats:%s", username)

	// Try to get from Redis cache first (30 minutes TTL)
	cachedJSON, err := services.RedisGetRawJSON(cacheKey)
	if err == nil && cachedJSON != "" {
		log.Printf("✅ Twitch stats cache HIT")
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, cachedJSON)
		return
	}
	log.Printf("🔍 Twitch stats cache MISS, fetching from API")

	stats, err := services.FetchTwitchStats(username)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Cache the stats for 30 minutes (1800 seconds)
	if err := services.RedisSetJSON(cacheKey, stats, 1800); err != nil {
		log.Printf("⚠️ Failed to cache Twitch stats: %v", err)
	} else {
		log.Printf("💾 Cached Twitch stats")
	}

	// Store stats in database
	database := db.GetDB()
	if database != nil {
		if err := db.InsertTwitchStats(database, username, stats.Followers, 0, false); err != nil {
			log.Printf("❌ Failed to save Twitch stats: %v", err)
		}
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
	cacheKey := fmt.Sprintf("leetcode:stats:%s", username)

	// Try to get from Redis cache first (30 minutes TTL)
	cachedJSON, err := services.RedisGetRawJSON(cacheKey)
	if err == nil && cachedJSON != "" {
		log.Printf("✅ LeetCode stats cache HIT")
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, cachedJSON)
		return
	}
	log.Printf("🔍 LeetCode stats cache MISS, fetching from API")

	stats, err := services.FetchLeetCodeStats(username)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// Cache the stats for 30 minutes (1800 seconds)
	if err := services.RedisSetJSON(cacheKey, stats, 1800); err != nil {
		log.Printf("⚠️ Failed to cache LeetCode stats: %v", err)
	} else {
		log.Printf("💾 Cached LeetCode stats")
	}

	// Store stats in database
	database := db.GetDB()
	if database != nil {
		if err := db.InsertLeetCodeStats(database, username, stats.SolvedCount, stats.Ranking, stats.Languages); err != nil {
			log.Printf("❌ Failed to save LeetCode stats: %v", err)
		}
	}

	c.JSON(http.StatusOK, stats)
}
