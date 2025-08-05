package models

type TwitchStats struct {
	DisplayName     string `json:"display_name"`
	Description     string `json:"description"`
	BroadcasterType string `json:"broadcaster_type"`
	Followers       int    `json:"followers"`
}
