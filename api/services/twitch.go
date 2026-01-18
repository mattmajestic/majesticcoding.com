package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"majesticcoding.com/api/models"
)

func getTwitchToken() (string, error) {
	cacheKey := "twitch:token:oauth"

	// Try to get from Redis cache first (1 hour TTL = 3600 seconds)
	cachedToken, err := RedisGet(cacheKey)
	if err == nil && cachedToken != "" {
		log.Printf("✅ Twitch token cache HIT")
		return cachedToken, nil
	}
	log.Printf("🔍 Twitch token cache MISS, fetching new token")

	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("TWITCH_CLIENT_ID or TWITCH_CLIENT_SECRET not set")
	}

	resp, err := http.PostForm("https://id.twitch.tv/oauth2/token", url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"grant_type":    {"client_credentials"},
	})
	if err != nil {
		return "", fmt.Errorf("failed to get twitch token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("twitch token API returned status %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	if result.AccessToken == "" {
		return "", fmt.Errorf("received empty access token from Twitch")
	}

	// Cache the token for 1 hour (3600 seconds) or use expires_in from response
	ttl := result.ExpiresIn
	if ttl == 0 {
		ttl = 3600 // fallback to 1 hour
	}

	if err := RedisSet(cacheKey, result.AccessToken, ttl); err != nil {
		log.Printf("⚠️ Failed to cache Twitch token: %v", err)
	} else {
		log.Printf("💾 Cached Twitch token for %d seconds", ttl)
	}

	return result.AccessToken, nil
}

func FetchTwitchStats(username string) (models.TwitchStats, error) {
	token, err := getTwitchToken()
	if err != nil {
		return models.TwitchStats{}, err
	}

	clientID := os.Getenv("TWITCH_CLIENT_ID")
	if clientID == "" {
		return models.TwitchStats{}, fmt.Errorf("TWITCH_CLIENT_ID not set")
	}

	url := fmt.Sprintf("https://api.twitch.tv/helix/users?login=%s", username)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.TwitchStats{}, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		return models.TwitchStats{}, fmt.Errorf("twitch API returned status %d", resp.StatusCode)
	}

	var result struct {
		Data []struct {
			DisplayName     string `json:"display_name"`
			Description     string `json:"description"`
			BroadcasterType string `json:"broadcaster_type"`
			ID              string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return models.TwitchStats{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Data) == 0 {
		return models.TwitchStats{}, fmt.Errorf("no user found for username: %s", username)
	}

	user := result.Data[0]

	// Fetch followers
	followersURL := fmt.Sprintf("https://api.twitch.tv/helix/channels/followers?broadcaster_id=%s", user.ID)
	req2, _ := http.NewRequest("GET", followersURL, nil)
	req2.Header.Set("Client-ID", clientID)
	req2.Header.Set("Authorization", "Bearer "+token)

	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		return models.TwitchStats{}, fmt.Errorf("followers request failed: %w", err)
	}
	defer resp2.Body.Close()

	// Check HTTP status code for followers endpoint
	if resp2.StatusCode != http.StatusOK {
		return models.TwitchStats{}, fmt.Errorf("followers API returned status %d", resp2.StatusCode)
	}

	var followResult struct {
		Total int `json:"total"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&followResult); err != nil {
		return models.TwitchStats{}, fmt.Errorf("failed to decode followers: %w", err)
	}

	return models.TwitchStats{
		DisplayName:     user.DisplayName,
		Description:     user.Description,
		BroadcasterType: user.BroadcasterType,
		Followers:       followResult.Total,
	}, nil
}
