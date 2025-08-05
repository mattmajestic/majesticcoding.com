// handlers/chat.go
package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/models"
)

func GetMessages(c *gin.Context) {
	Mu.Lock()
	defer Mu.Unlock()
	c.JSON(http.StatusOK, Messages)
}

func StartMessageCleanup() {
	go func() {
		for {
			time.Sleep(1 * time.Minute) // Check every minute
			cutoff := time.Now().Add(-10 * time.Minute)
			Mu.Lock()
			// Keep only messages newer than 10 minutes
			var filtered []models.Message
			for _, msg := range Messages {
				if msg.Timestamp.After(cutoff) {
					filtered = append(filtered, msg)
				}
			}
			Messages = filtered
			Mu.Unlock()
		}
	}()
}
