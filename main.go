package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := os.Getenv("DATABASE_URL")
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchYouTubeMetricsFromAPI(c *gin.Context) {
	resp, err := http.Get("https://mattmajestic.dev/youtube-metrics")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch YouTube metrics"})
		return
	}
	defer resp.Body.Close()
	c.DataFromReader(http.StatusOK, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

func startFFmpeg() {
    cmd := exec.Command("ffmpeg", "-f", "flv", "-i", "rtmp://localhost:1935/live", "-c:v", "libx264", "-f", "flv", "rtmp://localhost:1935/live/stream")
    if err := cmd.Start(); err != nil {
        log.Fatalf("Failed to start FFmpeg: %v", err)
    }
    log.Println("FFmpeg started to handle RTMP stream")
}

func main() {

	err_env := godotenv.Load()
	if err_env != nil {
		log.Fatalf("Error loading .env file: %v", err_env)
	}

	// Initialize Gin router
	router := gin.Default()

	_, err := db.Exec("INSERT INTO visits (visit_time) VALUES (NOW())")
	if err != nil {
		log.Printf("Error inserting visit time: %v", err)
	}

	// Serve static files
	router.GET("/styles.css", func(c *gin.Context) {
		c.File("./static/styles.css")
	})
	router.GET("/script.js", func(c *gin.Context) {
		c.File("./static/script.js")
	})
	router.GET("/clerk.js", func(c *gin.Context) {
		c.File("./static/clerk.js")
	})
	router.GET("/google-analytics.js", func(c *gin.Context) {
		c.File("./static/google-analytics.js")
	})
	router.GET("/download-resume-pdf", func(c *gin.Context) {
		c.File("./static/Matt Majestic Resume.pdf")
	})

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	// Define routes
	router.GET("/", func(c *gin.Context) {
		// Render the index template
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"CLERK_PUBLISHABLE_KEY": os.Getenv("CLERK_PUBLISHABLE_KEY"),
		})
	})

	router.GET("/login", func(c *gin.Context) {
		// Render the index template
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"CLERK_PUBLISHABLE_KEY": os.Getenv("CLERK_PUBLISHABLE_KEY"),
		})
	})

	// Route to fetch YouTube metrics from API
	router.GET("/youtube-metrics", fetchYouTubeMetricsFromAPI)

	// Route to serve the stream template
	router.GET("/stream", func(c *gin.Context) {
		tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "stream.tmpl"), filepath.Join("templates", "header.tmpl")))
		err := tmpl.ExecuteTemplate(c.Writer, "stream.tmpl", nil)
		if err != nil {
			c.String(http.StatusInternalServerError, "Template rendering error: %v", err)
		}
	})

	// Start FFmpeg to handle RTMP stream
	go startFFmpeg()

	// Run the server on port 8080
	router.Run(":8080")
}
