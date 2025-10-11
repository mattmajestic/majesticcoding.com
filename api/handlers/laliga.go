package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/services"
)

// GetLaLigaSchedule godoc
// @Summary Get La Liga weekly schedule
// @Description Returns matches for the next 7 days from La Liga
// @Tags La Liga
// @Success 200 {object} models.LaLigaScheduleResponse
// @Failure 500 {object} map[string]string
// @Router /laliga/schedule [get]
func GetLaLigaSchedule(c *gin.Context) {
	cacheKey := "laliga:schedule:weekly"

	// Try to get from Redis cache first (6 hours TTL = 21600 seconds)
	cachedJSON, err := services.RedisGetRawJSON(cacheKey)
	if err == nil && cachedJSON != "" {
		log.Printf("✅ La Liga schedule cache HIT")
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, cachedJSON)
		return
	}
	log.Printf("🔍 La Liga schedule cache MISS, fetching from API")

	matches, err := services.FetchLaLigaSchedule()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"matches": matches,
		"count":   len(matches),
	}

	// Cache the schedule for 6 hours (21600 seconds)
	if err := services.RedisSetJSON(cacheKey, response, 21600); err != nil {
		log.Printf("⚠️ Failed to cache La Liga schedule: %v", err)
	} else {
		log.Printf("💾 Cached La Liga schedule for 6 hours")
	}

	c.JSON(http.StatusOK, response)
}

// LaLigaWidget renders the La Liga widget template
func LaLigaWidget(c *gin.Context) {
	c.HTML(http.StatusOK, "laliga.tmpl", nil)
}
