// handlers/store.go
package handlers

import (
	"database/sql"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"majesticcoding.com/api/models"
	"majesticcoding.com/db"
)

var (
	Messages  []models.Message
	Clients   = make(map[*websocket.Conn]bool)
	Broadcast = make(chan models.Message)
	Mu        sync.Mutex
)

// LoadMessagesFromDB seeds the in-memory Messages slice from the database on startup
func LoadMessagesFromDB(database *sql.DB) {
	recent, err := db.GetRecentMessages(database, 50)
	if err != nil {
		log.Printf("⚠️  Could not load recent messages from DB: %v", err)
		return
	}
	Mu.Lock()
	Messages = recent
	Mu.Unlock()
	log.Printf("💬 Loaded %d recent messages from database", len(recent))
}
