package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Auth checks for the presence of an Authorization header.
// You can expand this to call Clerk’s verification API.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// TODO: Validate token with Clerk if needed

		c.Next()
	}
}
