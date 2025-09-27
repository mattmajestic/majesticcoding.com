package services

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"majesticcoding.com/api/models"
	"majesticcoding.com/db"

	"github.com/gempir/go-twitch-irc/v4"
)

var (
	messages     []models.TwitchMessage
	messagesLock sync.Mutex
	maxMessages  = 50
)

// StartTwitchChatFeed starts an anonymous Twitch client and stores messages
func StartTwitchChatFeed(channel string) {
	client := twitch.NewAnonymousClient()

	client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		messagesLock.Lock()
		defer messagesLock.Unlock()

		twitchMsg := models.TwitchMessage{
			Username:      msg.User.Name,
			DisplayName:   msg.User.DisplayName,
			Message:       msg.Message,
			Color:         msg.User.Color,
			Badges:        msg.User.Badges,
			IsMod:         msg.User.IsMod,
			IsVip:         msg.User.IsVip,
			IsBroadcaster: msg.User.IsBroadcaster,
			Time:          msg.Time,
		}

		// Store in database
		database := db.GetDB()
		if database != nil {
			if err := db.InsertTwitchMessage(database, twitchMsg); err != nil {
				log.Printf("❌ Failed to save Twitch message to database: %v", err)
			} else {
				log.Printf("💬 Saved Twitch message from %s", msg.User.DisplayName)
			}
		}

		// Check for !checkin command
		if strings.HasPrefix(strings.ToLower(msg.Message), "!checkin ") {
			location := strings.TrimSpace(msg.Message[9:]) // Remove "!checkin "
			if location != "" {
				log.Printf("🌍 Processing !checkin command from %s: %s", msg.User.DisplayName, location)
				go handleCheckinCommand(location, msg.User.DisplayName)
			}
		}

		// Keep in memory for quick access
		messages = append(messages, twitchMsg)
		if len(messages) > maxMessages {
			messages = messages[len(messages)-maxMessages:]
		}
	})

	client.OnConnect(func() {})
	client.Join(channel)

	go func() {
		if err := client.Connect(); err != nil {
			log.Printf("❌ Twitch connection failed: %v", err)
		}
	}()
}

// GetRecentMessages returns last Twitch messages
func GetRecentMessages() []models.TwitchMessage {
	messagesLock.Lock()
	defer messagesLock.Unlock()

	copied := make([]models.TwitchMessage, len(messages))
	copy(copied, messages)
	return copied
}

// handleCheckinCommand processes !checkin commands from Twitch chat
func handleCheckinCommand(location, username string) {
	// Call the geocode API endpoint with username
	geocodeURL := fmt.Sprintf("http://localhost:8080/api/geocode?city=%s&username=%s",
		url.QueryEscape(location), url.QueryEscape(username))

	resp, err := http.Get(geocodeURL)
	if err != nil {
		log.Printf("❌ Failed to call geocode API for %s: %v", location, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("✅ Successfully processed !checkin for %s from %s", location, username)
	} else {
		log.Printf("❌ Geocode API returned status %d for location: %s", resp.StatusCode, location)
	}
}
