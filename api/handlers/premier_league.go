package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/services"
)

// GetPremierLeagueSchedule godoc
// @Summary Get Premier League weekly schedule
// @Description Returns matches for the next 7 days from the Premier League
// @Tags Premier League
// @Success 200 {object} models.PLScheduleResponse
// @Failure 500 {object} map[string]string
// @Router /epl/schedule [get]
func GetPremierLeagueSchedule(c *gin.Context) {
	cacheKey := "epl:schedule:weekly"

	// Try to get from Redis cache first (7 days TTL = 604800 seconds)
	cachedJSON, err := services.RedisGetRawJSON(cacheKey)
	if err == nil && cachedJSON != "" {
		log.Printf("✅ EPL schedule cache HIT")
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, cachedJSON)
		return
	}
	log.Printf("🔍 EPL schedule cache MISS, fetching from API")

	matches, err := services.FetchPLSchedule()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"matches": matches,
		"count":   len(matches),
	}

	// Cache the schedule for 7 days (604800 seconds)
	if err := services.RedisSetJSON(cacheKey, response, 604800); err != nil {
		log.Printf("⚠️ Failed to cache EPL schedule: %v", err)
	} else {
		log.Printf("💾 Cached EPL schedule for 7 days")
	}

	c.JSON(http.StatusOK, response)
}

// EPLWidget renders the EPL widget template
func EPLWidget(c *gin.Context) {
	c.HTML(http.StatusOK, "epl.tmpl", nil)
}
