package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// sanitizeURL removes or shortens sensitive information from URLs
func sanitizeURL(url string) string {
	// Hide JWT tokens in WebSocket URLs
	if strings.Contains(url, "/ws/chat?token=") {
		parts := strings.Split(url, "?token=")
		if len(parts) > 1 {
			// Show only first 20 characters of token
			token := parts[1]
			if len(token) > 20 {
				token = token[:20] + "..."
			}
			return parts[0] + "?token=" + token
		}
	}
	return url
}

// CustomLogger returns a gin.HandlerFunc (middleware) that logs requests in a custom format
func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log after processing
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		// Build the URL with query parameters
		fullURL := path
		if raw != "" {
			fullURL = path + "?" + raw
		}

		// Sanitize sensitive information
		sanitizedURL := sanitizeURL(fullURL)

		// Format status code with color
		statusCodeStr := fmt.Sprintf("%d", statusCode)

		// Skip logging for very long WebSocket connections (> 5 seconds)
		// These are normal long-lived connections
		if strings.Contains(path, "/ws/") && latency > 5*time.Second {
			fmt.Printf("[GIN] %s | %s | WebSocket closed | %s | %s      %s\n",
				end.Format("2006/01/02 - 15:04:05"),
				statusCodeStr,
				clientIP,
				method,
				sanitizedURL)
		} else {
			fmt.Printf("[GIN] %s | %s | %v | %s | %s      %s\n",
				end.Format("2006/01/02 - 15:04:05"),
				statusCodeStr,
				latency,
				clientIP,
				method,
				sanitizedURL)
		}
	}
}
