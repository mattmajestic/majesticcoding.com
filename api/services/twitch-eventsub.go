package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"majesticcoding.com/api/models"
	"majesticcoding.com/db"
)

type EventSubClient struct {
	conn          *websocket.Conn
	sessionID     string
	isConnected   bool
	mu            sync.RWMutex
	reconnectChan chan bool
}

type EventSubMessage struct {
	Metadata struct {
		MessageID           string    `json:"message_id"`
		MessageType         string    `json:"message_type"`
		MessageTimestamp    time.Time `json:"message_timestamp"`
		SubscriptionType    string    `json:"subscription_type,omitempty"`
		SubscriptionVersion string    `json:"subscription_version,omitempty"`
	} `json:"metadata"`
	Payload struct {
		Session      *SessionPayload      `json:"session,omitempty"`
		Subscription *SubscriptionPayload `json:"subscription,omitempty"`
		Event        json.RawMessage      `json:"event,omitempty"`
	} `json:"payload"`
}

type SessionPayload struct {
	ID                      string    `json:"id"`
	Status                  string    `json:"status"`
	ConnectedAt             time.Time `json:"connected_at"`
	KeepaliveTimeoutSeconds int       `json:"keepalive_timeout_seconds"`
	ReconnectURL            string    `json:"reconnect_url"`
}

type SubscriptionPayload struct {
	ID        string                 `json:"id"`
	Status    string                 `json:"status"`
	Type      string                 `json:"type"`
	Version   string                 `json:"version"`
	Condition map[string]interface{} `json:"condition"`
	Transport struct {
		Method    string `json:"method"`
		SessionID string `json:"session_id"`
	} `json:"transport"`
	CreatedAt time.Time `json:"created_at"`
	Cost      int       `json:"cost"`
}

var eventSubClient *EventSubClient

// getTwitchUserToken - helper function to get current valid user token from database
func getTwitchUserToken() (string, error) {
	database := db.GetDB()
	if database == nil {
		return "", fmt.Errorf("database not available")
	}

	accessToken, _, _, _, expiresAt, err := db.GetTwitchToken(database)
	if err != nil {
		return "", err
	}

	// Check if token is expired
	if time.Now().After(expiresAt) {
		return "", fmt.Errorf("token expired at %v", expiresAt)
	}

	return accessToken, nil
}

func StartTwitchEventSub() error {
	eventSubClient = &EventSubClient{
		reconnectChan: make(chan bool, 1),
	}

	go eventSubClient.connect()
	go eventSubClient.handleReconnects()

	return nil
}

func (c *EventSubClient) connect() error {
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial("wss://eventsub.wss.twitch.tv/ws", nil)
	if err != nil {
		log.Printf("❌ Failed to connect to Twitch EventSub: %v", err)
		c.scheduleReconnect()
		return err
	}

	c.mu.Lock()
	c.conn = conn
	c.isConnected = true
	c.mu.Unlock()

	log.Println("🔗 Connected to Twitch EventSub WebSocket")

	defer func() {
		c.mu.Lock()
		c.isConnected = false
		c.conn.Close()
		c.mu.Unlock()
		c.scheduleReconnect()
	}()

	for {
		var message EventSubMessage
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("❌ Error reading EventSub message: %v", err)
			break
		}

		c.handleMessage(message)
	}

	return nil
}

func (c *EventSubClient) handleMessage(message EventSubMessage) {
	switch message.Metadata.MessageType {
	case "session_welcome":
		log.Println("📨 Received session welcome")
		if message.Payload.Session != nil {
			c.mu.Lock()
			c.sessionID = message.Payload.Session.ID
			c.mu.Unlock()
			log.Printf("🆔 Session ID: %s", c.sessionID)

			go c.createSubscriptions()
		}

	case "session_keepalive":
		log.Println("💓 Keepalive received")

	case "notification":
		log.Printf("🔔 Notification received: %s", message.Metadata.SubscriptionType)
		c.handleNotification(message)

	case "session_reconnect":
		log.Println("🔄 Session reconnect requested")
		if message.Payload.Session != nil && message.Payload.Session.ReconnectURL != "" {
			go c.reconnectTo(message.Payload.Session.ReconnectURL)
		}

	case "revocation":
		log.Printf("❌ Subscription revoked: %s", message.Metadata.SubscriptionType)

	default:
		log.Printf("❓ Unknown message type: %s", message.Metadata.MessageType)
	}
}

func (c *EventSubClient) handleNotification(message EventSubMessage) {
	database := db.GetDB()
	if database == nil {
		log.Println("❌ Database not available for EventSub notification")
		return
	}

	switch message.Metadata.SubscriptionType {
	case "channel.follow":
		var followEvent struct {
			UserID     string    `json:"user_id"`
			UserLogin  string    `json:"user_login"`
			UserName   string    `json:"user_name"`
			FollowedAt time.Time `json:"followed_at"`
		}

		if err := json.Unmarshal(message.Payload.Event, &followEvent); err != nil {
			log.Printf("❌ Error parsing follow event: %v", err)
			return
		}

		follower := models.TwitchFollower{
			UserID:     followEvent.UserID,
			UserLogin:  followEvent.UserLogin,
			UserName:   followEvent.UserName,
			FollowedAt: followEvent.FollowedAt,
		}

		if err := db.InsertTwitchFollower(database, follower); err != nil {
			log.Printf("❌ Failed to save follower: %v", err)
		} else {
			log.Printf("✅ New follower saved: %s", follower.UserName)
		}

	case "channel.raid":
		var raidEvent struct {
			FromBroadcasterUserID    string `json:"from_broadcaster_user_id"`
			FromBroadcasterUserLogin string `json:"from_broadcaster_user_login"`
			FromBroadcasterUserName  string `json:"from_broadcaster_user_name"`
			ToBroadcasterUserID      string `json:"to_broadcaster_user_id"`
			ToBroadcasterUserLogin   string `json:"to_broadcaster_user_login"`
			ToBroadcasterUserName    string `json:"to_broadcaster_user_name"`
			Viewers                  int    `json:"viewers"`
		}

		if err := json.Unmarshal(message.Payload.Event, &raidEvent); err != nil {
			log.Printf("❌ Error parsing raid event: %v", err)
			return
		}

		raid := models.TwitchRaid{
			FromBroadcasterUserID:    raidEvent.FromBroadcasterUserID,
			FromBroadcasterUserLogin: raidEvent.FromBroadcasterUserLogin,
			FromBroadcasterUserName:  raidEvent.FromBroadcasterUserName,
			ToBroadcasterUserID:      raidEvent.ToBroadcasterUserID,
			ToBroadcasterUserLogin:   raidEvent.ToBroadcasterUserLogin,
			ToBroadcasterUserName:    raidEvent.ToBroadcasterUserName,
			Viewers:                  raidEvent.Viewers,
		}

		if err := db.InsertTwitchRaid(database, raid); err != nil {
			log.Printf("❌ Failed to save raid: %v", err)
		} else {
			log.Printf("✅ New raid saved: %s -> %s (%d viewers)", raid.FromBroadcasterUserName, raid.ToBroadcasterUserName, raid.Viewers)
		}

	case "channel.subscribe":
		var subEvent struct {
			UserID               string `json:"user_id"`
			UserLogin            string `json:"user_login"`
			UserName             string `json:"user_name"`
			BroadcasterUserID    string `json:"broadcaster_user_id"`
			BroadcasterUserLogin string `json:"broadcaster_user_login"`
			BroadcasterUserName  string `json:"broadcaster_user_name"`
			Tier                 string `json:"tier"`
			IsGift               bool   `json:"is_gift"`
		}

		if err := json.Unmarshal(message.Payload.Event, &subEvent); err != nil {
			log.Printf("❌ Error parsing sub event: %v", err)
			return
		}

		sub := models.TwitchSub{
			UserID:               subEvent.UserID,
			UserLogin:            subEvent.UserLogin,
			UserName:             subEvent.UserName,
			BroadcasterUserID:    subEvent.BroadcasterUserID,
			BroadcasterUserLogin: subEvent.BroadcasterUserLogin,
			BroadcasterUserName:  subEvent.BroadcasterUserName,
			Tier:                 subEvent.Tier,
			IsGift:               subEvent.IsGift,
		}

		if err := db.InsertTwitchSub(database, sub); err != nil {
			log.Printf("❌ Failed to save subscription: %v", err)
		} else {
			log.Printf("✅ New subscription saved: %s (Tier %s)", sub.UserName, sub.Tier)
		}

	case "channel.cheer":
		var cheerEvent struct {
			UserID               string `json:"user_id"`
			UserLogin            string `json:"user_login"`
			UserName             string `json:"user_name"`
			BroadcasterUserID    string `json:"broadcaster_user_id"`
			BroadcasterUserLogin string `json:"broadcaster_user_login"`
			BroadcasterUserName  string `json:"broadcaster_user_name"`
			IsAnonymous          bool   `json:"is_anonymous"`
			Message              string `json:"message"`
			Bits                 int    `json:"bits"`
		}

		if err := json.Unmarshal(message.Payload.Event, &cheerEvent); err != nil {
			log.Printf("❌ Error parsing cheer event: %v", err)
			return
		}

		bits := models.TwitchBits{
			UserID:               cheerEvent.UserID,
			UserLogin:            cheerEvent.UserLogin,
			UserName:             cheerEvent.UserName,
			BroadcasterUserID:    cheerEvent.BroadcasterUserID,
			BroadcasterUserLogin: cheerEvent.BroadcasterUserLogin,
			BroadcasterUserName:  cheerEvent.BroadcasterUserName,
			IsAnonymous:          cheerEvent.IsAnonymous,
			Message:              &cheerEvent.Message,
			Bits:                 cheerEvent.Bits,
		}

		if err := db.InsertTwitchBits(database, bits); err != nil {
			log.Printf("❌ Failed to save bits: %v", err)
		} else {
			log.Printf("✅ New bits saved: %s cheered %d bits", bits.UserName, bits.Bits)
		}
	}
}

func (c *EventSubClient) createSubscriptions() {
	c.mu.RLock()
	sessionID := c.sessionID
	c.mu.RUnlock()

	if sessionID == "" {
		log.Println("❌ No session ID available for creating subscriptions")
		return
	}

	// Get broadcaster user ID using stored user token
	log.Println("🔍 Getting broadcaster user ID for majesticcodingtwitch...")

	// Try to get stored user access token (required for EventSub WebSocket)
	userToken, err := getTwitchUserToken()
	if err != nil {
		log.Printf("❌ No valid user access token found: %v", err)
		log.Println("💡 Visit https://majesticcoding.com/api/twitch/oauth/start to authenticate")
		return
	}

	log.Println("✅ Using stored user access token for EventSub")

	clientID := os.Getenv("TWITCH_CLIENT_ID")
	url := "https://api.twitch.tv/helix/users?login=majesticcodingtwitch"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", "Bearer "+userToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("❌ Failed to get user info: %v", err)
		return
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || len(result.Data) == 0 {
		log.Printf("❌ Failed to decode user data")
		return
	}

	broadcasterUserID := result.Data[0].ID
	log.Printf("✅ Found broadcaster user ID: %s", broadcasterUserID)

	// Start with RAIDS - easiest to test!
	log.Println("🎯 Creating RAIDS subscription...")
	condition := map[string]interface{}{
		"to_broadcaster_user_id": broadcasterUserID,
	}

	if err := c.createSubscription("channel.raid", "1", condition, sessionID, userToken); err != nil {
		log.Printf("❌ Failed to create channel.raid subscription: %v", err)
	} else {
		log.Printf("✅ Created channel.raid subscription! 🎉")
		log.Println("💡 Raid events will be captured when someone raids your channel")
	}
}

func (c *EventSubClient) createSubscription(subscriptionType, version string, condition map[string]interface{}, sessionID, accessToken string) error {
	clientID := os.Getenv("TWITCH_CLIENT_ID")

	if clientID == "" {
		return fmt.Errorf("TWITCH_CLIENT_ID not set")
	}

	if accessToken == "" {
		return fmt.Errorf("access token is required")
	}

	payload := map[string]interface{}{
		"type":      subscriptionType,
		"version":   version,
		"condition": condition,
		"transport": map[string]interface{}{
			"method":     "websocket",
			"session_id": sessionID,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.twitch.tv/helix/eventsub/subscriptions", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		// Read the response body for more details
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("subscription request failed with status: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *EventSubClient) scheduleReconnect() {
	select {
	case c.reconnectChan <- true:
	default:
	}
}

func (c *EventSubClient) handleReconnects() {
	for range c.reconnectChan {
		log.Println("🔄 Attempting to reconnect in 5 seconds...")
		time.Sleep(5 * time.Second)
		go c.connect()
	}
}

func (c *EventSubClient) reconnectTo(url string) {
	log.Printf("🔄 Reconnecting to: %s", url)

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Printf("❌ Failed to reconnect: %v", err)
		c.scheduleReconnect()
		return
	}

	c.mu.Lock()
	if c.conn != nil {
		c.conn.Close()
	}
	c.conn = conn
	c.isConnected = true
	c.mu.Unlock()

	log.Println("✅ Successfully reconnected to Twitch EventSub")

	for {
		var message EventSubMessage
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("❌ Error reading EventSub message after reconnect: %v", err)
			break
		}

		c.handleMessage(message)
	}
}
