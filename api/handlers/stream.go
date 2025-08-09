package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/models"
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
	awsURL := os.Getenv("AWS_STREAMING_URL")
	stream := models.NewStream("", awsURL)
	if stream.IsActive {
		c.String(http.StatusOK, "true")
	} else {
		c.String(http.StatusOK, "false")
	}
}
