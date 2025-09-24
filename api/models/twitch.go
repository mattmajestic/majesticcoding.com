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

type TwitchFollower struct {
	ID         int       `json:"id"`
	UserID     string    `json:"user_id"`
	UserLogin  string    `json:"user_login"`
	UserName   string    `json:"user_name"`
	FollowedAt time.Time `json:"followed_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type TwitchRaid struct {
	ID                       int       `json:"id"`
	FromBroadcasterUserID    string    `json:"from_broadcaster_user_id"`
	FromBroadcasterUserLogin string    `json:"from_broadcaster_user_login"`
	FromBroadcasterUserName  string    `json:"from_broadcaster_user_name"`
	ToBroadcasterUserID      string    `json:"to_broadcaster_user_id"`
	ToBroadcasterUserLogin   string    `json:"to_broadcaster_user_login"`
	ToBroadcasterUserName    string    `json:"to_broadcaster_user_name"`
	Viewers                  int       `json:"viewers"`
	CreatedAt                time.Time `json:"created_at"`
}

type TwitchSub struct {
	ID                   int       `json:"id"`
	UserID               string    `json:"user_id"`
	UserLogin            string    `json:"user_login"`
	UserName             string    `json:"user_name"`
	BroadcasterUserID    string    `json:"broadcaster_user_id"`
	BroadcasterUserLogin string    `json:"broadcaster_user_login"`
	BroadcasterUserName  string    `json:"broadcaster_user_name"`
	Tier                 string    `json:"tier"`
	IsGift               bool      `json:"is_gift"`
	GifterUserID         *string   `json:"gifter_user_id,omitempty"`
	GifterUserLogin      *string   `json:"gifter_user_login,omitempty"`
	GifterUserName       *string   `json:"gifter_user_name,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
}

type TwitchBits struct {
	ID                   int       `json:"id"`
	UserID               string    `json:"user_id"`
	UserLogin            string    `json:"user_login"`
	UserName             string    `json:"user_name"`
	BroadcasterUserID    string    `json:"broadcaster_user_id"`
	BroadcasterUserLogin string    `json:"broadcaster_user_login"`
	BroadcasterUserName  string    `json:"broadcaster_user_name"`
	IsAnonymous          bool      `json:"is_anonymous"`
	Message              *string   `json:"message,omitempty"`
	Bits                 int       `json:"bits"`
	CreatedAt            time.Time `json:"created_at"`
}
