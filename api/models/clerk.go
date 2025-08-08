package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClerkSession struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func StoreSession(c *gin.Context) {
	var session ClerkSession
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// Store session in memory or set a cookie
	c.SetCookie("clerk_user", session.ID, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "session stored"})
}
