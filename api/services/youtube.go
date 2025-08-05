package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func FetchYouTubeStats() (map[string]interface{}, error) {
	apiKey := os.Getenv("YT_API_KEY")
	channelID := os.Getenv("YT_CHANNEL_ID")

	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/channels?part=snippet,statistics&id=%s&key=%s", channelID, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Items []struct {
			Snippet struct {
				Title string `json:"title"`
			} `json:"snippet"`
			Statistics struct {
				SubscriberCount string `json:"subscriberCount"`
				ViewCount       string `json:"viewCount"`
				VideoCount      string `json:"videoCount"`
			} `json:"statistics"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("no items returned")
	}

	item := result.Items[0]

	// Return simplified map
	return map[string]interface{}{
		"channelTitle": item.Snippet.Title,
		"subscribers":  item.Statistics.SubscriberCount,
		"views":        item.Statistics.ViewCount,
		"videos":       item.Statistics.VideoCount,
	}, nil
}
