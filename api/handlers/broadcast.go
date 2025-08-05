package handlers

import (
	"log"
)

func StartBroadcaster() {
	go func() {
		for msg := range Broadcast {
			Mu.Lock()
			for client := range Clients {
				if err := client.WriteJSON(msg); err != nil {
					log.Println("Broadcast error:", err)
					client.Close()
					delete(Clients, client)
				}
			}
			Mu.Unlock()
		}
	}()
}
