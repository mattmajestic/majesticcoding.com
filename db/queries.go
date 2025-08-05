package db

import (
	"database/sql"

	"majesticcoding.com/api/models"
)

// InsertMessage inserts a new message into the messages table
func InsertMessage(db *sql.DB, content string) error {
	_, err := db.Exec(`INSERT INTO messages (content) VALUES ($1)`, content)
	return err
}

// GetRecentMessages fetches the most recent messages up to the given limit
func GetRecentMessages(db *sql.DB, limit int) ([]models.Message, error) {
	rows, err := db.Query(`SELECT id, content, created_at FROM messages ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.Username, &msg.Content, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
