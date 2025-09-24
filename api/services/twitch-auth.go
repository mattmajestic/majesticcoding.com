package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type TwitchTokenResponse struct {
	AccessToken string   `json:"access_token"`
	ExpiresIn   int      `json:"expires_in"`
	TokenType   string   `json:"token_type"`
	Scope       []string `json:"scope"`
}

// GetTwitchAppAccessToken gets an app access token using client credentials flow
// This is suitable for EventSub subscriptions that don't require user-specific permissions
func GetTwitchAppAccessToken() (string, error) {
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("TWITCH_CLIENT_ID or TWITCH_CLIENT_SECRET not set")
	}

	// Prepare form data for client credentials flow
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "client_credentials")
	// Add scopes needed for EventSub subscriptions
	data.Set("scope", "moderator:read:followers channel:read:subscriptions bits:read")

	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	var tokenResponse TwitchTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	log.Printf("✅ Got Twitch app access token (expires in %d seconds)", tokenResponse.ExpiresIn)
	return tokenResponse.AccessToken, nil
}

// GetTwitchUserAccessToken gets a user access token (requires user authorization)
// This is more complex and requires implementing OAuth flow with redirect URIs
// For now, we'll focus on app access tokens for EventSub
func GetTwitchUserAccessToken(code, redirectURI string) (string, error) {
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("TWITCH_CLIENT_ID or TWITCH_CLIENT_SECRET not set")
	}

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	var tokenResponse TwitchTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}
