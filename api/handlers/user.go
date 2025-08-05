package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthStatusHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusOK, gin.H{"loggedIn": false})
		return
	}

	// Send token to Clerk's /v1/me endpoint
	req, err := http.NewRequest("GET", "https://api.clerk.dev/v1/me", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "request creation failed"})
		return
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		c.JSON(http.StatusOK, gin.H{"loggedIn": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"loggedIn": true})
}
