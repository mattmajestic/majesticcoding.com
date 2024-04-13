package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// FetchYouTubeMetricsFromAPI fetches YouTube metrics from an API
func FetchYouTubeMetricsFromAPI(c *gin.Context) {
    resp, err := http.Get("https://mattmajestic.dev/youtube-metrics")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch YouTube metrics"})
        return
    }
    defer resp.Body.Close()
    c.DataFromReader(http.StatusOK, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

func ExampleFunction() {
    // Example function in utils.go
}
