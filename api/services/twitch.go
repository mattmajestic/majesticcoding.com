package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"majesticcoding.com/api/models"
)

func getTwitchToken() (string, error) {
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
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
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
