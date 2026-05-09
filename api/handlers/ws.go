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
	Subprotocols: []string{"clerk-auth", "supabase-auth"},
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

// getUsernameFromAuth extracts a display name from a Clerk JWT
func getUsernameFromAuth(r *http.Request) string {
	tokenString := getTokenFromRequest(r)
	if tokenString == "" {
		return generateAnonUsername()
	}

	claims, err := verifyClerkToken(r.Context(), tokenString)
	if err != nil {
		log.Printf("Clerk token verify failed: %v", err)
		return generateAnonUsername()
	}

	if username, ok := claims["username"].(string); ok && username != "" {
		return "✓ " + username
	}
	if name, ok := claims["name"].(string); ok && name != "" {
		return "✓ " + name
	}
	if email, ok := claims["email"].(string); ok && email != "" {
		return "✓ " + email
	}

	// JWT verified but no profile data — show user ID prefix so at least it's not anon
	log.Printf("Clerk user verified (sub=%s) but no name/email returned", claims["sub"])
	if sub, ok := claims["sub"].(string); ok && len(sub) > 8 {
		return "✓ " + sub[:12]
	}
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
