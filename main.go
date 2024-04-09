package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func fetchYouTubeMetricsFromAPI(c *gin.Context) {
	resp, err := http.Get("https://mattmajestic.dev/youtube-metrics")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch YouTube metrics"})
		return
	}
	defer resp.Body.Close()

	c.DataFromReader(http.StatusOK, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", nil)
	})

	router.GET("/youtube-metrics", fetchYouTubeMetricsFromAPI)

	router.Run(":8080")
}
