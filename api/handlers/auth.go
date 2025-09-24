package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/models"
)

// AuthStatus returns the current user's authentication status
func AuthStatus(c *gin.Context) {
	// For now, just return a simple status
	// This can be enhanced with Supabase token verification if needed
	c.JSON(http.StatusOK, models.AuthResponse{
		Success: false,
		Message: "Please use Supabase authentication",
	})
}

// AuthStatusHandler function already exists in user.go
