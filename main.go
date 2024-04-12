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
    // Initialize Gin router
    router := gin.Default()

    // Serve static files
    router.GET("/styles.css", func(c *gin.Context) {
        c.File("./static/styles.css")
    })
    router.GET("/script.js", func(c *gin.Context) {
        c.File("./static/script.js")
    })

    // Load HTML templates
    router.LoadHTMLGlob("templates/*")

    // Define routes
    router.GET("/", func(c *gin.Context) {
        // Render the index template
        c.HTML(http.StatusOK, "index.tmpl", nil)
    })

    // Route to fetch YouTube metrics from API
    router.GET("/youtube-metrics", fetchYouTubeMetricsFromAPI)

    // Run the server on port 8080
    router.Run(":8080")
}
