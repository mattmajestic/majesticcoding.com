package services

import (
	"log"
	"sync"

	"majesticcoding.com/api/models"

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

		messages = append(messages, models.TwitchMessage{
			Username:      msg.User.Name,
			DisplayName:   msg.User.DisplayName,
			Message:       msg.Message,
			Color:         msg.User.Color,
			Badges:        msg.User.Badges,
			IsMod:         msg.User.IsMod,
			IsVip:         msg.User.IsVip,
			IsBroadcaster: msg.User.IsBroadcaster,
			Time:          msg.Time,
		})

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
