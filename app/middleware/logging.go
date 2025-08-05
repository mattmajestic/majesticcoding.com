package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger logs each request with method, path, status, and latency.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.Printf("[%d] %s %s (%s)", status, c.Request.Method, c.Request.URL.Path, latency)
	}
}
