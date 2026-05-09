package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"majesticcoding.com/api/models"
)

const twitchTokenCacheKey = "twitch:token:oauth"

func invalidateTwitchTokenCache() {
	RedisDelete(twitchTokenCacheKey)
}

func getTwitchToken() (string, error) {
	// Try to get from Redis cache first (1 hour TTL = 3600 seconds)
	cachedToken, err := RedisGet(twitchTokenCacheKey)
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

	if err := RedisSet(twitchTokenCacheKey, result.AccessToken, ttl); err != nil {
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

	// 401: cached token expired — clear and retry once with fresh token
	if resp.StatusCode == http.StatusUnauthorized {
		invalidateTwitchTokenCache()
		return FetchTwitchStats(username)
	}
	// 403: log the full response body so we can see Twitch's reason
	if resp.StatusCode == http.StatusForbidden {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("⚠️ Twitch 403 body: %s | Client-ID=%s", string(body), clientID)
		return models.TwitchStats{}, fmt.Errorf("twitch 403: %s", string(body))
	}
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

	// Fetch followers — requires moderator:read:followers user scope; fall back to 0 if unavailable
	followers := 0
	followersURL := fmt.Sprintf("https://api.twitch.tv/helix/channels/followers?broadcaster_id=%s", user.ID)
	req2, _ := http.NewRequest("GET", followersURL, nil)
	req2.Header.Set("Client-ID", clientID)
	req2.Header.Set("Authorization", "Bearer "+token)

	if resp2, err := http.DefaultClient.Do(req2); err == nil {
		defer resp2.Body.Close()
		if resp2.StatusCode == http.StatusOK {
			var followResult struct {
				Total int `json:"total"`
			}
			if err := json.NewDecoder(resp2.Body).Decode(&followResult); err == nil {
				followers = followResult.Total
			}
		}
	}

	return models.TwitchStats{
		DisplayName:     user.DisplayName,
		Description:     user.Description,
		BroadcasterType: user.BroadcasterType,
		Followers:       followers,
	}, nil
}
