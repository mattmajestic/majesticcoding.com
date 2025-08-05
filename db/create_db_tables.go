package db

import (
	"database/sql"
)

func CreateTables(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS example (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL
	);`)
	return err
}

func CreateMessagesTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			content TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return err
}
