package models

import "time"

type Message struct {
	Content     string
	Username    string
	Timestamp   time.Time
	DisplayTime string
}
