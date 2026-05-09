package handlers

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	clerkuser "github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gin-gonic/gin"
)

// clerkFrontendDomain decodes the Clerk domain from a publishable key.
// Format: pk_[live|test]_<base64(domain + "$")>
func clerkFrontendDomain(publishableKey string) string {
	encoded := strings.TrimPrefix(strings.TrimPrefix(publishableKey, "pk_live_"), "pk_test_")
	decoded, err := base64.StdEncoding.DecodeString(encoded + "=") // pad if needed
	if err != nil {
		// try without padding
		decoded, err = base64.RawStdEncoding.DecodeString(encoded)
		if err != nil {
			return ""
		}
	}
	return strings.TrimSuffix(string(decoded), "$")
}

// InitClerk sets the Clerk secret key — call this after loading env vars
func InitClerk() {
	key := os.Getenv("CLERK_SECRET_KEY")
	pk := os.Getenv("CLERK_PUBLISHABLE_KEY")
	log.Printf("🔑 Clerk init: publishable=%s... secret=%s...", pk[:min(20, len(pk))], key[:min(15, len(key))])
	clerk.SetKey(key)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// verifyClerkToken verifies a Clerk JWT and fetches user details
func verifyClerkToken(ctx context.Context, tokenString string) (map[string]interface{}, error) {
	claims, err := jwt.Verify(ctx, &jwt.VerifyParams{Token: tokenString})
	if err != nil {
		return nil, err
	}

	raw := map[string]interface{}{"sub": claims.Subject}

	user, err := clerkuser.Get(ctx, claims.Subject)
	if err == nil {
		if user.Username != nil && *user.Username != "" {
			raw["username"] = *user.Username
		}
		name := ""
		if user.FirstName != nil {
			name = *user.FirstName
		}
		if user.LastName != nil && *user.LastName != "" {
			name = strings.TrimSpace(name + " " + *user.LastName)
		}
		if name != "" {
			raw["name"] = name
		}
		if len(user.EmailAddresses) > 0 {
			raw["email"] = user.EmailAddresses[0].EmailAddress
		}
	}

	return raw, nil
}

// ClerkAuthMiddleware verifies Clerk session tokens from Authorization header
func ClerkAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		token := strings.TrimSpace(authHeader[7:])
		claims, err := jwt.Verify(c.Request.Context(), &jwt.VerifyParams{Token: token})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("user_id", claims.Subject)
		c.Next()
	}
}

// ClerkConfigHandler returns the Clerk publishable key for the frontend
func ClerkConfigHandler(c *gin.Context) {
	publishableKey := os.Getenv("CLERK_PUBLISHABLE_KEY")
	if publishableKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Clerk configuration missing"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"publishableKey": publishableKey})
}

// ClerkAuthStatusHandler returns whether the current request is authenticated
func ClerkAuthStatusHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
		return
	}

	token := strings.TrimSpace(authHeader[7:])
	claims, err := jwt.Verify(c.Request.Context(), &jwt.VerifyParams{Token: token})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"authenticated": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"authenticated": true, "user_id": claims.Subject})
}

// RenderClerkAuth renders the Clerk sign-in page
func RenderClerkAuth(c *gin.Context) {
	pk := os.Getenv("CLERK_PUBLISHABLE_KEY")
	c.HTML(http.StatusOK, "clerk-auth.tmpl", gin.H{
		"title":          "Login - Majestic Coding",
		"publishableKey": pk,
		"clerkDomain":    clerkFrontendDomain(pk),
	})
}
