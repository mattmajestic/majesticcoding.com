package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gin-gonic/gin"
)

var clerkClient, _ = clerk.NewClient(os.Getenv("CLERK_SECRET_KEY"))

func AuthStatus(c *gin.Context) {
	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"logged_in": false})
		return
	}

	session, err := clerkClient.VerifyToken(token)
	if err != nil || session == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"logged_in": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logged_in": true,
	})
}
