// handlers/chat.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMessages godoc
// @Summary Get chat messages
// @Description Returns all chat messages
// @Tags Chat
// @Success 200 {array} models.Message
// @Router /chat [get]
func GetMessages(c *gin.Context) {
	Mu.Lock()
	defer Mu.Unlock()
	c.JSON(http.StatusOK, Messages)
}
