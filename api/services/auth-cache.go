package services

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"time"

	"majesticcoding.com/db"
)

const (
	// Cache TTL - 15 minutes
	SessionTTL = 15 * time.Minute
)

// hashToken creates a secure hash of the JWT token for storage
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", hash)
}

// GetCachedUserData retrieves user data from cache, returns nil if not found or expired
func GetCachedUserData(database *sql.DB, token string) (map[string]interface{}, error) {
	if database == nil {
		return nil, fmt.Errorf("database not available")
	}

	tokenHash := hashToken(token)
	userData, err := db.GetCachedSession(database, tokenHash)
	if err != nil {
		// Cache miss or expired - not an error, just return nil
		return nil, nil
	}

	return userData, nil
}

// SetCachedUserData stores user data in cache with TTL
func SetCachedUserData(database *sql.DB, token string, userData map[string]interface{}) error {
	if database == nil {
		return fmt.Errorf("database not available")
	}

	tokenHash := hashToken(token)
	expiresAt := time.Now().Add(SessionTTL)

	err := db.SetCachedSession(database, tokenHash, userData, expiresAt)
	if err != nil {
		log.Printf("Failed to cache session: %v", err)
		return err
	}

	return nil
}

// InvalidateUserSession removes a user's session from cache
func InvalidateUserSession(database *sql.DB, token string) error {
	if database == nil {
		return fmt.Errorf("database not available")
	}

	tokenHash := hashToken(token)
	err := db.InvalidateSession(database, tokenHash)
	if err != nil {
		log.Printf("Failed to invalidate session: %v", err)
		return err
	}

	return nil
}

// StartSessionCleanup starts a goroutine to periodically clean up expired sessions
func StartSessionCleanup(database *sql.DB) {
	go func() {
		ticker := time.NewTicker(5 * time.Minute) // Cleanup every 5 minutes
		defer ticker.Stop()

		for range ticker.C {
			if database != nil {
				err := db.CleanupExpiredSessions(database)
				if err != nil {
					log.Printf("Failed to cleanup expired sessions: %v", err)
				} else {
					log.Printf("Successfully cleaned up expired auth sessions")
				}
			}
		}
	}()
}
