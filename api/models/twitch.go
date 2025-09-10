package models

import "time"

type TwitchStats struct {
	DisplayName     string `json:"display_name"`
	Description     string `json:"description"`
	BroadcasterType string `json:"broadcaster_type"`
	Followers       int    `json:"followers"`
}

type TwitchMessage struct {
	ID            int            `json:"id"`
	Username      string         `json:"username"`
	DisplayName   string         `json:"display_name"`
	Message       string         `json:"message"`
	Color         string         `json:"color,omitempty"`
	Badges        map[string]int `json:"badges,omitempty"`
	IsMod         bool           `json:"is_mod"`
	IsVip         bool           `json:"is_vip"`
	IsBroadcaster bool           `json:"is_broadcaster"`
	Time          time.Time      `json:"time"`
	CreatedAt     time.Time      `json:"created_at"`
}
