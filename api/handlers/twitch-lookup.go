package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/services"
)

// TwitchUserLookupHandler looks up Twitch user ID by username
// @Summary Look up Twitch user ID
// @Description Get numeric user ID for a Twitch username
// @Tags twitch
// @Produce json
// @Param username query string true "Twitch username to lookup"
// @Success 200 {object} map[string]interface{}
// @Router /api/twitch/lookup [get]
func TwitchUserLookupHandler(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username parameter is required"})
		return
	}

	userID, err := services.GetTwitchUserID(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"user_id":  userID,
	})
}
