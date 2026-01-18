package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"majesticcoding.com/api/services"
)

type speechWSStart struct {
	Type        string `json:"type"`
	ContentType string `json:"contentType"`
	Filename    string `json:"filename"`
	Language    string `json:"language"`
}

type speechWSResponse struct {
	Type    string `json:"type"`
	Text    string `json:"text,omitempty"`
	IsFinal bool   `json:"isFinal,omitempty"`
	Error   string `json:"error,omitempty"`
}

type speechWSState struct {
	mu            sync.Mutex
	contentType   string
	filename      string
	language      string
	data          []byte
	lastAudioAt   time.Time
	transcribing  bool
	lastTranscript string
}

const maxSpeechBufferBytes = 8 * 1024 * 1024

func SpeechWebSocket(c *gin.Context) {
	token := getSupabaseTokenFromRequest(c.Request)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth token"})
		return
	}
	if _, err := verifySupabaseToken(token); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	state := &speechWSState{
		filename: "speech.webm",
		language: "en-US",
	}

	done := make(chan struct{})
	go speechTranscriptionLoop(conn, state, done)

	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			close(done)
			return
		}

		switch messageType {
		case websocket.TextMessage:
			var start speechWSStart
			if err := json.Unmarshal(payload, &start); err != nil {
				continue
			}
			if strings.ToLower(start.Type) == "start" {
				state.mu.Lock()
				if start.ContentType != "" {
					state.contentType = start.ContentType
				}
				if start.Filename != "" {
					state.filename = start.Filename
				}
				if start.Language != "" {
					state.language = start.Language
				}
				state.mu.Unlock()
			}
			if strings.ToLower(start.Type) == "stop" {
				sendFinalTranscript(conn, state)
				conn.WriteJSON(speechWSResponse{Type: "done"})
				close(done)
				return
			}
		case websocket.BinaryMessage:
			state.mu.Lock()
			if len(state.data)+len(payload) > maxSpeechBufferBytes {
				state.mu.Unlock()
				conn.WriteJSON(speechWSResponse{Type: "error", Error: "audio buffer too large"})
				close(done)
				return
			}
			state.data = append(state.data, payload...)
			state.lastAudioAt = time.Now()
			state.mu.Unlock()
		}
	}
}

func speechTranscriptionLoop(conn *websocket.Conn, state *speechWSState, done <-chan struct{}) {
	ticker := time.NewTicker(800 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			state.mu.Lock()
			if state.transcribing || len(state.data) == 0 {
				state.mu.Unlock()
				continue
			}
			if time.Since(state.lastAudioAt) < 300*time.Millisecond {
				state.mu.Unlock()
				continue
			}
			state.transcribing = true
			snapshot := append([]byte(nil), state.data...)
			contentType := state.contentType
			filename := state.filename
			state.mu.Unlock()

			text, err := services.TranscribeAudio(bytes.NewReader(snapshot), filename, contentType)
			state.mu.Lock()
			state.transcribing = false
			if err != nil {
				state.mu.Unlock()
				conn.WriteJSON(speechWSResponse{Type: "error", Error: err.Error()})
				continue
			}
			text = strings.TrimSpace(text)
			if text != "" && text != state.lastTranscript {
				state.lastTranscript = text
				state.mu.Unlock()
				conn.WriteJSON(speechWSResponse{Type: "transcript", Text: text, IsFinal: false})
				continue
			}
			state.mu.Unlock()
		}
	}
}

func sendFinalTranscript(conn *websocket.Conn, state *speechWSState) {
	state.mu.Lock()
	if len(state.data) == 0 || state.transcribing {
		state.mu.Unlock()
		return
	}
	snapshot := append([]byte(nil), state.data...)
	contentType := state.contentType
	filename := state.filename
	state.mu.Unlock()

	text, err := services.TranscribeAudio(bytes.NewReader(snapshot), filename, contentType)
	if err != nil {
		conn.WriteJSON(speechWSResponse{Type: "error", Error: err.Error()})
		return
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	conn.WriteJSON(speechWSResponse{Type: "transcript", Text: text, IsFinal: true})
}

// token parsing helper lives in ws_auth.go
