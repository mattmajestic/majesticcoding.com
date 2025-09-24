package db

import (
	"database/sql"
	"encoding/json"
	"time"
)

// GetCachedSession retrieves cached session data by token hash
func GetCachedSession(db *sql.DB, tokenHash string) (map[string]interface{}, error) {
	query := `
		SELECT user_data
		FROM bronze.auth_sessions
		WHERE token_hash = $1 AND expires_at > CURRENT_TIMESTAMP
	`

	var userDataJSON []byte
	err := db.QueryRow(query, tokenHash).Scan(&userDataJSON)
	if err != nil {
		return nil, err
	}

	var userData map[string]interface{}
	if err := json.Unmarshal(userDataJSON, &userData); err != nil {
		return nil, err
	}

	return userData, nil
}

// SetCachedSession stores session data with expiration
func SetCachedSession(db *sql.DB, tokenHash string, userData map[string]interface{}, expiresAt time.Time) error {
	userDataJSON, err := json.Marshal(userData)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO bronze.auth_sessions (token_hash, user_data, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (token_hash) DO UPDATE SET
			user_data = EXCLUDED.user_data,
			expires_at = EXCLUDED.expires_at
	`

	_, err = db.Exec(query, tokenHash, userDataJSON, expiresAt)
	return err
}

// InvalidateSession removes a cached session
func InvalidateSession(db *sql.DB, tokenHash string) error {
	query := `DELETE FROM bronze.auth_sessions WHERE token_hash = $1`
	_, err := db.Exec(query, tokenHash)
	return err
}

// CleanupExpiredSessions removes expired sessions (call periodically)
func CleanupExpiredSessions(db *sql.DB) error {
	query := `DELETE FROM bronze.auth_sessions WHERE expires_at < CURRENT_TIMESTAMP`
	_, err := db.Exec(query)
	return err
}
