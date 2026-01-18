package handlers

import (
	"net/http"
	"strings"
)

func getSupabaseTokenFromRequest(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		}
		return strings.TrimSpace(authHeader)
	}

	protocolHeader := r.Header.Get("Sec-WebSocket-Protocol")
	if protocolHeader == "" {
		return ""
	}

	parts := strings.Split(protocolHeader, ",")
	var token string
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		if value == "supabase-auth" {
			continue
		}
		token = value
		break
	}

	if token == "" {
		return ""
	}
	return token
}
