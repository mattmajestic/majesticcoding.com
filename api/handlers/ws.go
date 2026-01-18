package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
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
		return isAllowedWSOrigin(r)
	},
	Subprotocols: []string{"supabase-auth"},
}

func generateAnonUsername() string {
	return fmt.Sprintf(namesgenerator.GetRandomName(0)+"_%02d", time.Now().UnixNano()%10000)
}

func isAllowedWSOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}

	allowedOrigins := strings.Split(os.Getenv("WS_ALLOWED_ORIGINS"), ",")
	if len(allowedOrigins) > 0 && strings.TrimSpace(allowedOrigins[0]) != "" {
		for _, allowed := range allowedOrigins {
			if strings.TrimSpace(allowed) == origin {
				return true
			}
		}
		return false
	}

	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}

	return parsed.Host == r.Host
}

// getUsernameFromAuth attempts to get username from Supabase auth token
func getUsernameFromAuth(r *http.Request) string {
	tokenString := getSupabaseTokenFromRequest(r)
	if tokenString == "" {
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

	log.Printf("✅ User %s connected to chat", username)

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
	// Simple in-memory count of connected clients
	Mu.Lock()
	count := len(Clients)
	Mu.Unlock()

	log.Printf("✅ Connected chat users: %d", count)
	c.JSON(http.StatusOK, gin.H{"user_count": count, "source": "memory"})
}
