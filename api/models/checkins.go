package models

import "time"

type Checkin struct {
	ID          int       `json:"id"`
	Username    string    `json:"username,omitempty"`
	Lat         float64   `json:"lat"`
	Lon         float64   `json:"lon"`
	City        string    `json:"city,omitempty"`
	Country     string    `json:"country,omitempty"`
	CheckinTime time.Time `json:"checkin_time"`
}
