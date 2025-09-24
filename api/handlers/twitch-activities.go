package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/db"
)

// TwitchFollowersHandler returns recent Twitch followers
// @Summary Get recent Twitch followers
// @Description Returns a list of recent Twitch followers
// @Tags twitch
// @Produce json
// @Param limit query int false "Number of followers to return (default: 10)"
// @Success 200 {array} models.TwitchFollower
// @Router /api/twitch/followers [get]
func TwitchFollowersHandler(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not available"})
		return
	}

	followers, err := db.GetRecentTwitchFollowers(database, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch followers"})
		return
	}

	c.JSON(http.StatusOK, followers)
}

// TwitchRaidsHandler returns recent Twitch raids
// @Summary Get recent Twitch raids
// @Description Returns a list of recent Twitch raids
// @Tags twitch
// @Produce json
// @Param limit query int false "Number of raids to return (default: 10)"
// @Success 200 {array} models.TwitchRaid
// @Router /api/twitch/raids [get]
func TwitchRaidsHandler(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not available"})
		return
	}

	raids, err := db.GetRecentTwitchRaids(database, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch raids"})
		return
	}

	c.JSON(http.StatusOK, raids)
}

// TwitchSubsHandler returns recent Twitch subscriptions
// @Summary Get recent Twitch subscriptions
// @Description Returns a list of recent Twitch subscriptions
// @Tags twitch
// @Produce json
// @Param limit query int false "Number of subscriptions to return (default: 10)"
// @Success 200 {array} models.TwitchSub
// @Router /api/twitch/subs [get]
func TwitchSubsHandler(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not available"})
		return
	}

	subs, err := db.GetRecentTwitchSubs(database, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscriptions"})
		return
	}

	c.JSON(http.StatusOK, subs)
}

// TwitchBitsHandler returns recent Twitch bits/cheers
// @Summary Get recent Twitch bits
// @Description Returns a list of recent Twitch bits/cheers
// @Tags twitch
// @Produce json
// @Param limit query int false "Number of bits to return (default: 10)"
// @Success 200 {array} models.TwitchBits
// @Router /api/twitch/bits [get]
func TwitchBitsHandler(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not available"})
		return
	}

	bits, err := db.GetRecentTwitchBits(database, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bits"})
		return
	}

	c.JSON(http.StatusOK, bits)
}
