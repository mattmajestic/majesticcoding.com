package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/services"
)

func TwitchMessagesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, services.GetRecentMessages())
}
