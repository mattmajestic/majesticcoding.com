package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/db"
)

var (
	twitchOauthState = "twitch-majestic-state"
)

type TwitchTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

type SavedTwitchToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
	Scopes       string    `json:"scopes"`
}

// Save token to database
func saveTwitchToken(token *TwitchTokenResponse) error {
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("database not available")
	}

	expiresAt := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	return db.SaveTwitchToken(database, token.AccessToken, token.RefreshToken, token.TokenType, token.Scope, expiresAt)
}

// Load token from database
func loadTwitchToken() (*SavedTwitchToken, error) {
	database := db.GetDB()
	if database == nil {
		return nil, fmt.Errorf("database not available")
	}

	accessToken, refreshToken, tokenType, scopes, expiresAt, err := db.GetTwitchToken(database)
	if err != nil {
		return nil, err
	}

	token := &SavedTwitchToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    tokenType,
		ExpiresAt:    expiresAt,
		Scopes:       scopes,
	}

	// Check if token is expired and try to refresh
	if time.Now().After(token.ExpiresAt) {
		log.Println("Saved Twitch token expired, attempting to refresh...")

		if token.RefreshToken == "" {
			return nil, fmt.Errorf("token expired and no refresh token available")
		}

		// Try to refresh the token
		refreshedToken, err := refreshTwitchToken(token.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %v", err)
		}

		// Create new saved token with refreshed data
		newSavedToken := SavedTwitchToken{
			AccessToken:  refreshedToken.AccessToken,
			RefreshToken: refreshedToken.RefreshToken,
			ExpiresAt:    time.Now().Add(time.Duration(refreshedToken.ExpiresIn) * time.Second),
			TokenType:    refreshedToken.TokenType,
			Scopes:       refreshedToken.Scope,
		}

		// Save the refreshed token to database
		if err := db.UpdateTwitchToken(database, newSavedToken.AccessToken, newSavedToken.RefreshToken, newSavedToken.ExpiresAt); err == nil {
			log.Println("Successfully refreshed and saved new Twitch token to database")
		} else {
			log.Printf("Failed to save refreshed Twitch token to database: %v", err)
		}

		return &newSavedToken, nil
	}

	return token, nil
}

// Refresh access token using refresh token
func refreshTwitchToken(refreshTokenStr string) (*TwitchTokenResponse, error) {
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("missing client credentials")
	}

	// Prepare form data
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshTokenStr)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	// Create request
	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("token refresh failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TwitchTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %v", err)
	}

	// If no new refresh token provided, keep the old one
	if tokenResp.RefreshToken == "" {
		tokenResp.RefreshToken = refreshTokenStr
	}

	return &tokenResp, nil
}

// Exchange authorization code for token
func exchangeTwitchCodeForToken(code, redirectURI string) (*TwitchTokenResponse, error) {
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("missing client credentials")
	}

	// Prepare form data
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURI)

	// Create request
	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	log.Printf("Twitch token response status: %d", resp.StatusCode)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TwitchTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %v", err)
	}

	return &tokenResp, nil
}

// GET /api/twitch/oauth/start
func TwitchOAuthHandler(c *gin.Context) {
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	if clientID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "TWITCH_CLIENT_ID not set"})
		return
	}

	redirectURI := "https://majesticcoding.com/api/twitch/oauth/callback"
	scopes := "moderator:read:followers channel:read:subscriptions bits:read"

	authURL := fmt.Sprintf(
		"https://id.twitch.tv/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s",
		clientID,
		url.QueryEscape(redirectURI),
		url.QueryEscape(scopes),
		twitchOauthState,
	)

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GET /api/twitch/oauth/callback
func TwitchOAuthCallbackHandler(c *gin.Context) {
	if c.Query("state") != twitchOauthState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "state mismatch"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing authorization code"})
		return
	}

	log.Printf("Attempting Twitch token exchange with code: %s...", code[:10])

	redirectURI := "https://majesticcoding.com/api/twitch/oauth/callback"
	tokenResp, err := exchangeTwitchCodeForToken(code, redirectURI)
	if err != nil {
		log.Printf("TWITCH EXCHANGE ERROR: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   "token exchange failed",
			"details": err.Error(),
		})
		return
	}

	// Save token to database
	if err := saveTwitchToken(tokenResp); err != nil {
		log.Printf("Failed to save Twitch token: %v", err)
	} else {
		log.Println("Twitch token saved successfully")
	}

	log.Printf("Twitch auth successful! Token expires in: %d seconds", tokenResp.ExpiresIn)
	log.Printf("Scopes: %s", tokenResp.Scope)

	c.JSON(http.StatusOK, gin.H{
		"ok":         true,
		"message":    "Twitch connected successfully! EventSub should now work.",
		"expires_in": tokenResp.ExpiresIn,
		"token_type": tokenResp.TokenType,
		"scopes":     tokenResp.Scope,
	})
}

// GET /api/twitch/status - check if Twitch user token is available
func TwitchStatusHandler(c *gin.Context) {
	token, err := loadTwitchToken()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"connected": false,
			"message":   "Not connected. Visit /api/twitch/oauth/start to authenticate for EventSub.",
			"error":     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"connected":  true,
		"message":    "Twitch user token available for EventSub!",
		"expires_at": token.ExpiresAt,
		"scopes":     token.Scopes,
	})
}

// InitTwitchClient - call this on startup to check for saved tokens
func InitTwitchClient() {
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Println("WARNING: TWITCH_CLIENT_ID and TWITCH_CLIENT_SECRET environment variables not set")
		log.Println("Twitch EventSub will not be available")
		return
	}

	// Try to load saved token
	if savedToken, err := loadTwitchToken(); err == nil {
		log.Printf("✅ Twitch user token loaded from database (expires: %v)", savedToken.ExpiresAt)
		log.Printf("📋 Scopes: %s", savedToken.Scopes)
		log.Println("🎉 EventSub WebSocket ready!")
	} else {
		log.Printf("No valid saved Twitch token found: %v", err)
		log.Printf("Twitch ready - ClientID: %s...", clientID[:10])
		log.Println("💡 Visit /api/twitch/oauth/start to authenticate for EventSub")
	}
}

// GetTwitchUserToken - helper function for EventSub to get current valid token
func GetTwitchUserToken() (string, error) {
	token, err := loadTwitchToken()
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}
