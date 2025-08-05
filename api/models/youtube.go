package models

type YouTubeStats struct {
	ChannelName string `json:"channel_name"`
	Subscribers int    `json:"subscribers"`
	Views       int    `json:"views"`
	Videos      int    `json:"videos"`
}
