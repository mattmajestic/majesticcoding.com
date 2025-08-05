package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	goaway "github.com/TwiN/go-away"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/moby/moby/pkg/namesgenerator"
	"majesticcoding.com/api/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func generateAnonUsername() string {
	return fmt.Sprintf(namesgenerator.GetRandomName(0)+"_%02d", time.Now().UnixNano()%10000)
}

func ChatWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	username := generateAnonUsername()

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

		Mu.Lock()
		Messages = append(Messages, msg)
		Mu.Unlock()

		Broadcast <- msg
	}
}
