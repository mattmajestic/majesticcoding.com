package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetClerkSession(c *gin.Context) {
	userID, err := c.Cookie("clerk_user")
	if err != nil || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"logged_in": false})
		return
	}

	// If you stored username/email in another cookie or session store, fetch them here if needed

	c.JSON(http.StatusOK, gin.H{
		"logged_in": true,
		"user_id":   userID,
	})
}
