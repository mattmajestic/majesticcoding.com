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

	resp, err := http.PostForm("https://id.twitch.tv/oauth2/token", url.Values{
		"client_id":     {os.Getenv("TWITCH_CLIENT_ID")},
		"client_secret": {os.Getenv("TWITCH_CLIENT_SECRET")},
		"grant_type":    {"client_credentials"},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
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
	url := fmt.Sprintf("https://api.twitch.tv/helix/users?login=%s", username)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.TwitchStats{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			DisplayName     string `json:"display_name"`
			Description     string `json:"description"`
			BroadcasterType string `json:"broadcaster_type"`
			ID              string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || len(result.Data) == 0 {
		return models.TwitchStats{}, fmt.Errorf("failed to fetch user data")
	}

	user := result.Data[0]

	// Fetch followers
	followersURL := fmt.Sprintf("https://api.twitch.tv/helix/channels/followers?broadcaster_id=%s", user.ID)
	req2, _ := http.NewRequest("GET", followersURL, nil)
	req2.Header.Set("Client-ID", clientID)
	req2.Header.Set("Authorization", "Bearer "+token)

	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		return models.TwitchStats{}, err
	}
	defer resp2.Body.Close()

	var followResult struct {
		Total int `json:"total"`
	}
	if err := json.NewDecoder(resp2.Body).Decode(&followResult); err != nil {
		return models.TwitchStats{}, err
	}

	return models.TwitchStats{
		DisplayName:     user.DisplayName,
		Description:     user.Description,
		BroadcasterType: user.BroadcasterType,
		Followers:       followResult.Total,
	}, nil
}
