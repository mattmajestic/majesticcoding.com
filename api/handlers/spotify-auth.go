package handlers

import (
	"context"
	"encoding/base64"
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
	spotify "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"majesticcoding.com/db"
)

var (
	spAuth     *spotifyauth.Authenticator
	oauthState = "majestic-state" // replace with a random string if you want CSRF defense
	spClient   *spotify.Client    // demo: single-user in-memory client
)

type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type SavedToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// Removed tokenFile constant - now using database

// Custom transport to add Bearer token to all requests
type tokenTransport struct {
	token string
	base  http.RoundTripper
}

func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(req)
}

// Save token to database
func saveToken(token *SpotifyTokenResponse) error {
	database := db.GetDB()
	if database == nil {
		return fmt.Errorf("database not available")
	}
	
	expiresAt := time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	return db.SaveSpotifyToken(database, token.AccessToken, token.RefreshToken, token.TokenType, expiresAt)
}

// Load token from database
func loadToken() (*SavedToken, error) {
	database := db.GetDB()
	if database == nil {
		return nil, fmt.Errorf("database not available")
	}

	accessToken, refreshToken, tokenType, expiresAt, err := db.GetSpotifyToken(database)
	if err != nil {
		return nil, err
	}

	token := &SavedToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    tokenType,
		ExpiresAt:    expiresAt,
	}

	// Check if token is expired and try to refresh
	if time.Now().After(token.ExpiresAt) {
		log.Println("Saved token expired, attempting to refresh...")
		
		if token.RefreshToken == "" {
			return nil, fmt.Errorf("token expired and no refresh token available")
		}
		
		// Try to refresh the token
		refreshedToken, err := refreshSpotifyToken(token.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %v", err)
		}
		
		// Create new saved token with refreshed data
		newSavedToken := SavedToken{
			AccessToken:  refreshedToken.AccessToken,
			RefreshToken: refreshedToken.RefreshToken,
			ExpiresAt:    time.Now().Add(time.Duration(refreshedToken.ExpiresIn) * time.Second),
			TokenType:    refreshedToken.TokenType,
		}
		
		// Save the refreshed token to database
		database := db.GetDB()
		if database != nil {
			if err := db.UpdateSpotifyToken(database, newSavedToken.AccessToken, newSavedToken.RefreshToken, newSavedToken.ExpiresAt); err == nil {
				log.Println("Successfully refreshed and saved new token to database")
			} else {
				log.Printf("Failed to save refreshed token to database: %v", err)
			}
		}
		
		return &newSavedToken, nil
	}

	return token, nil
}

// Refresh access token using refresh token
func refreshSpotifyToken(refreshTokenStr string) (*SpotifyTokenResponse, error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("missing client credentials")
	}

	// Create basic auth header
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	
	// Prepare form data
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshTokenStr)
	
	// Create request
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	
	req.Header.Set("Authorization", "Basic "+auth)
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
	
	var tokenResp SpotifyTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %v", err)
	}
	
	// If no new refresh token provided, keep the old one
	if tokenResp.RefreshToken == "" {
		tokenResp.RefreshToken = refreshTokenStr
	}
	
	return &tokenResp, nil
}

// InitSpotifyClient - call this on startup to initialize client if possible
func InitSpotifyClient() {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	redirectURI := os.Getenv("SPOTIFY_REDIRECT_URI")

	if clientID == "" || clientSecret == "" {
		log.Println("WARNING: SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET environment variables not set")
		log.Println("Spotify client will not be available")
		return
	}

	if redirectURI == "" {
		log.Println("WARNING: SPOTIFY_REDIRECT_URI environment variable not set")
		log.Println("Spotify client will not be available")
		return
	}

	// Initialize the authenticator after env vars are loaded
	spAuth = spotifyauth.New(
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithClientSecret(clientSecret),
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying),
	)

	// Try to load saved token
	if savedToken, err := loadToken(); err == nil {
		log.Printf("Loading saved Spotify token (expires: %v)", savedToken.ExpiresAt.Format("2006-01-02 15:04:05"))

		// Create HTTP client with saved token
		httpClient := &http.Client{
			Transport: &tokenTransport{
				token: savedToken.AccessToken,
				base:  http.DefaultTransport,
			},
		}

		spClient = spotify.New(httpClient)
		log.Println("Spotify client ready with saved token!")
	} else {
		log.Printf("No valid saved token found: %v", err)
		log.Printf("Spotify ready - ClientID: %s...", clientID[:10])
		log.Printf("Redirect URI: %s", redirectURI)
		log.Println("Visit /api/spotify/login to authenticate for user data access")
	}
}

// Manual token exchange using direct HTTP request
func exchangeCodeForToken(code, redirectURI string) (*SpotifyTokenResponse, error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("missing client credentials")
	}

	// Create basic auth header
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	// Prepare form data
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	// Create request
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Basic "+auth)
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

	log.Printf("Spotify token response status: %d", resp.StatusCode)
	log.Printf("Spotify token response body: %s", string(body))

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp SpotifyTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %v", err)
	}

	return &tokenResp, nil
}

// GET /api/spotify/login
func SpotifyLogin(c *gin.Context) {
	if spAuth == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Spotify authentication not configured",
		})
		return
	}
	url := spAuth.AuthURL(oauthState) + "&show_dialog=true"
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GET /callback
func SpotifyCallback(c *gin.Context) {
	if spAuth == nil {
		log.Println("ERROR: spAuth is nil - Spotify not initialized")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Spotify authentication not configured",
		})
		return
	}

	if c.Query("state") != oauthState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "state mismatch"})
		return
	}
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
		return
	}

	log.Printf("Attempting token exchange with code: %s...", code[:10])

	// Use manual token exchange
	redirectURI := os.Getenv("SPOTIFY_REDIRECT_URI")
	tokenResp, err := exchangeCodeForToken(code, redirectURI)
	if err != nil {
		log.Printf("SPOTIFY EXCHANGE ERROR: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   "token exchange failed",
			"details": err.Error(),
		})
		return
	}

	// Create HTTP client with Bearer token transport
	httpClient := &http.Client{
		Transport: &tokenTransport{
			token: tokenResp.AccessToken,
			base:  http.DefaultTransport,
		},
	}

	spClient = spotify.New(httpClient)

	// Save token to file for persistence
	if err := saveToken(tokenResp); err != nil {
		log.Printf("Failed to save token: %v", err)
	} else {
		log.Println("Token saved successfully")
	}

	// Log token info for debugging (don't log actual tokens in production!)
	log.Printf("Spotify auth successful! Token expires in: %d seconds", tokenResp.ExpiresIn)
	log.Printf("Refresh token available: %v", tokenResp.RefreshToken != "")

	c.JSON(http.StatusOK, gin.H{
		"ok":         true,
		"message":    "Spotify connected successfully!",
		"expires_in": tokenResp.ExpiresIn,
		"token_type": tokenResp.TokenType,
	})
}

// GET /api/spotify/status - check if Spotify is connected
func SpotifyStatus(c *gin.Context) {
	if spClient == nil {
		c.JSON(http.StatusOK, gin.H{
			"connected": false,
			"message":   "Not connected. Visit /api/spotify/login to authenticate.",
		})
		return
	}

	// Try to get current user to verify connection
	ctx := context.Background()
	user, err := spClient.CurrentUser(ctx)
	if err != nil {
		log.Printf("Spotify client error: %v", err)
		spClient = nil // Reset client if it's not working
		c.JSON(http.StatusOK, gin.H{
			"connected": false,
			"message":   "Connection expired. Visit /api/spotify/login to re-authenticate.",
			"error":     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"connected": true,
		"user":      user.DisplayName,
		"user_id":   user.ID,
		"message":   "Connected to Spotify!",
	})
}
