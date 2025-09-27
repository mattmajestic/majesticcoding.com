package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/models"
	"majesticcoding.com/api/services"
)

func StreamHandler(c *gin.Context) {
	awsURL := os.Getenv("AWS_STREAMING_URL")
	stream := models.NewStream("", awsURL)
	c.HTML(http.StatusOK, "live.tmpl", gin.H{
		"IsStreaming": stream.IsActive,
		"StreamURL":   stream.URL,
	})
}

// StreamStatusHandler godoc
// @Summary Stream Status
// @Description Returns whether the stream is currently active
// @Tags Stream
// @Success 200 {string} string "true or false"
// @Router /stream/status [get]
func StreamStatusHandler(c *gin.Context) {
	cacheKey := "stream:status:3m"

	// Try to get from Redis cache first (3 minutes TTL = 180 seconds)
	cachedStatus, err := services.RedisGet(cacheKey)
	if err == nil && cachedStatus != "" {
		log.Printf("✅ Stream status cache HIT: %s", cachedStatus)
		c.String(http.StatusOK, cachedStatus)
		return
	}
	log.Printf("🔍 Stream status cache MISS, checking AWS IVS")

	awsURL := os.Getenv("AWS_STREAMING_URL")
	stream := models.NewStream("", awsURL)

	var status string
	if stream.IsActive {
		status = "true"
	} else {
		status = "false"
	}

	// Cache the status for 3 minutes (180 seconds)
	if err := services.RedisSet(cacheKey, status, 180); err != nil {
		log.Printf("⚠️ Failed to cache stream status: %v", err)
	} else {
		log.Printf("💾 Cached stream status '%s' for 3 minutes", status)
	}

	c.String(http.StatusOK, status)
}
