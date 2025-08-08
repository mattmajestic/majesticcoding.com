package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthStatus(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"logged_in": false})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"logged_in": false})
		return
	}

	req, _ := http.NewRequest("GET", "https://api.clerk.dev/v1/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Clerk-Version", "v1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		c.JSON(http.StatusUnauthorized, gin.H{"logged_in": false})
		return
	}
	defer resp.Body.Close()

	var user map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&user)

	c.JSON(http.StatusOK, gin.H{
		"logged_in": true,
		"user":      user,
	})
}
