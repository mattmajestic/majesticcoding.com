// handlers/store.go
package handlers

import (
	"sync"

	"github.com/gorilla/websocket"
	"majesticcoding.com/api/models"
)

var (
	Messages  []models.Message
	Clients   = make(map[*websocket.Conn]bool)
	Broadcast = make(chan models.Message)
	Mu        sync.Mutex
)
