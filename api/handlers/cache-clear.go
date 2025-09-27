package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/services"
)

// ClearStatsCache clears all stats cache entries
func ClearStatsCache(c *gin.Context) {
	patterns := []string{
		"github:stats:*",
		"youtube:stats:*",
		"twitch:stats:*",
		"leetcode:stats:*",
		"epl:schedule:*",
		"laliga:schedule:*",
		"checkins:recent:*",
		"geocode:*",
	}

	cleared := 0
	for _, pattern := range patterns {
		if err := services.RedisClearPattern(pattern); err != nil {
			log.Printf("Failed to clear pattern %s: %v", pattern, err)
		} else {
			cleared++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  fmt.Sprintf("Cleared cache for %d patterns", cleared),
		"patterns": patterns,
	})
}
