package handlers

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/models"
)

var clerkClient, _ = clerk.NewClient(os.Getenv("CLERK_SECRET_KEY"))

// AuthStatus returns the current user's authentication status and user data
func AuthStatus(c *gin.Context) {
	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if token == "" {
		c.JSON(http.StatusOK, models.AuthResponse{
			Success: false,
			Message: "No authorization token provided",
		})
		return
	}

	session, err := clerkClient.VerifyToken(token)
	if err != nil || session == nil {
		c.JSON(http.StatusOK, models.AuthResponse{
			Success: false,
			Message: "Invalid or expired token",
		})
		return
	}

	userID := session.Subject
	clerkUser, err := clerkClient.Users().Read(userID)
	if err != nil || clerkUser == nil {
		c.JSON(http.StatusOK, models.AuthResponse{
			Success: false,
			Message: "User not found",
		})
		return
	}

	// Extract email from Clerk user
	email := ""
	if len(clerkUser.EmailAddresses) > 0 {
		email = clerkUser.EmailAddresses[0].EmailAddress
	}

	// Create our user model
	user := models.User{
		ID:        clerkUser.ID,
		Email:     email,
		Username:  clerkUser.Username,
		FirstName: clerkUser.FirstName,
		LastName:  clerkUser.LastName,
		CreatedAt: time.Unix(clerkUser.CreatedAt/1000, 0), // Clerk returns timestamp in milliseconds
		UpdatedAt: time.Unix(clerkUser.UpdatedAt/1000, 0),
	}

	userSession := models.UserSession{
		User:      user,
		LoggedIn:  true,
		SessionID: session.ID,
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Success: true,
		Message: "User authenticated successfully",
		User:    &userSession,
	})
}

// LoginHandler handles the login process after Clerk authentication
func LoginHandler(c *gin.Context) {
	// This endpoint is called after successful Clerk authentication
	// The frontend will call this with the JWT token to establish server-side session
	AuthStatus(c)
}

// LogoutHandler handles user logout
func LogoutHandler(c *gin.Context) {
	// Clear any server-side session data if needed
	c.JSON(http.StatusOK, models.AuthResponse{
		Success: true,
		Message: "Logged out successfully",
	})
}
