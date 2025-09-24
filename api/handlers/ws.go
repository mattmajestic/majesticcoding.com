package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	goaway "github.com/TwiN/go-away"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/moby/moby/pkg/namesgenerator"
	"majesticcoding.com/api/models"
	"majesticcoding.com/db"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func generateAnonUsername() string {
	return fmt.Sprintf(namesgenerator.GetRandomName(0)+"_%02d", time.Now().UnixNano()%10000)
}

// getUsernameFromAuth attempts to get username from Supabase auth token
func getUsernameFromAuth(r *http.Request) string {
	// Check for Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// Check for token in query parameters (common for WebSockets)
		token := r.URL.Query().Get("token")
		if token != "" {
			authHeader = "Bearer " + token
		}
	}

	if authHeader == "" {
		return generateAnonUsername()
	}

	// Extract token
	tokenString := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = strings.TrimSpace(authHeader[7:])
	} else {
		return generateAnonUsername()
	}

	// Verify Supabase token
	user, err := verifySupabaseToken(tokenString)
	if err != nil {
		log.Printf("Auth verification failed for chat: %v", err)
		return generateAnonUsername()
	}

	// Try to get username from user metadata first
	if userMetadata, ok := user["user_metadata"].(map[string]interface{}); ok {
		// Check for various username fields from different providers
		if username, ok := userMetadata["user_name"].(string); ok && username != "" {
			return "✓ " + username // GitHub, Twitch username
		}
		if username, ok := userMetadata["preferred_username"].(string); ok && username != "" {
			return "✓ " + username // Some OAuth providers
		}
		if username, ok := userMetadata["name"].(string); ok && username != "" {
			return "✓ " + username // Display name
		}
	}

	// Fallback to email if no username found
	if email, ok := user["email"].(string); ok && email != "" {
		return "✓ " + email
	}

	// Fallback to anonymous
	return generateAnonUsername()
}

func ChatWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	username := getUsernameFromAuth(c.Request)

	Mu.Lock()
	Clients[conn] = true
	Mu.Unlock()

	// Send existing messages to the new client
	Mu.Lock()
	for _, msg := range Messages {
		if err := conn.WriteJSON(msg); err != nil {
			log.Println("Send error:", err)
		}
	}
	Mu.Unlock()

	for {
		var msg models.Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("Read error:", err)
			Mu.Lock()
			delete(Clients, conn)
			Mu.Unlock()
			break
		}

		msg.Content = goaway.Censor(msg.Content)

		msg.Username = username
		msg.Timestamp = time.Now()
		msg.DisplayTime = msg.Timestamp.Format("15:04:05")

		// Store message in database
		database := db.GetDB()
		if database != nil {
			if err := db.InsertChatMessage(database, username, msg.Content); err != nil {
				log.Printf("❌ Failed to save chat message to database: %v", err)
			} else {
				log.Printf("💬 Saved chat message from %s: %s", username, msg.Content)
			}
		}

		Mu.Lock()
		Messages = append(Messages, msg)
		Mu.Unlock()

		Broadcast <- msg
	}
}

func ChatUserCount(c *gin.Context) {
	Mu.Lock()
	count := len(Clients)
	Mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"user_count": count})
}
