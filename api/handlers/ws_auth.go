package handlers

import (
	"log"
	"net/http"
	"strings"
)

func getTokenFromRequest(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		}
		return strings.TrimSpace(authHeader)
	}

	protocolHeader := r.Header.Get("Sec-WebSocket-Protocol")
	if protocolHeader == "" {
		log.Println("🔑 WS auth: no Authorization or Sec-WebSocket-Protocol header")
		return ""
	}

	// known protocol names (not tokens) — skip them
	knownProtocols := map[string]bool{"clerk-auth": true, "supabase-auth": true}

	parts := strings.Split(protocolHeader, ",")
	var token string
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" || knownProtocols[value] {
			continue
		}
		token = value
		break
	}

	if token == "" {
		log.Printf("🔑 WS auth: protocol header present but no token found: %q", protocolHeader)
		return ""
	}
	tokenPreview := token
	if len(tokenPreview) > 20 {
		tokenPreview = tokenPreview[:20]
	}
	log.Printf("🔑 WS auth: extracted token (first 20 chars): %s...", tokenPreview)
	return token
}
