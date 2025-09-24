package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type TwitchUsersResponse struct {
	Data []struct {
		ID          string `json:"id"`
		Login       string `json:"login"`
		DisplayName string `json:"display_name"`
	} `json:"data"`
}

// GetTwitchUserID fetches the numeric user ID for a given username
func GetTwitchUserID(username string) (string, error) {
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	accessToken := os.Getenv("TWITCH_ACCESS_TOKEN")

	// If no access token is set, try to get one using client credentials
	if accessToken == "" {
		token, err := GetTwitchAppAccessToken()
		if err != nil {
			return "", fmt.Errorf("failed to get app access token: %v", err)
		}
		accessToken = token
	}

	if clientID == "" {
		return "", fmt.Errorf("TWITCH_CLIENT_ID not set")
	}

	url := fmt.Sprintf("https://api.twitch.tv/helix/users?login=%s", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Client-ID", clientID)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response TwitchUsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if len(response.Data) == 0 {
		return "", fmt.Errorf("user not found: %s", username)
	}

	return response.Data[0].ID, nil
}
