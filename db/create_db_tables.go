package db

import (
	"database/sql"
	"fmt"
	"log"
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
	// First, add username column if it doesn't exist
	_, err := db.Exec(`
		ALTER TABLE messages ADD COLUMN IF NOT EXISTS username VARCHAR(50)
	`)
	if err != nil {
		log.Printf("Note: Could not add username column (may already exist): %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50),
			content TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return err
}

func CreateCheckinsTable(db *sql.DB) error {
	// First, drop the existing table if it has the wrong column types
	_, err := db.Exec(`DROP TABLE IF EXISTS checkins;`)
	if err != nil {
		fmt.Printf("WARNING: Failed to drop existing checkins table: %v\n", err)
	}

	// Create table with correct floating point types
	_, err = db.Exec(`
		CREATE TABLE checkins (
			id SERIAL PRIMARY KEY,
			lat DOUBLE PRECISION NOT NULL,
			lon DOUBLE PRECISION NOT NULL,
			city TEXT,
			country TEXT,
			checkin_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		fmt.Printf("ERROR: Failed to create checkins table: %v\n", err)
		return err
	}

	fmt.Printf("SUCCESS: Checkins table created with correct DOUBLE PRECISION columns\n")
	return nil
}

func CreateSpotifyTokensTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS spotify_tokens (
			id SERIAL PRIMARY KEY,
			access_token TEXT NOT NULL,
			refresh_token TEXT,
			token_type TEXT DEFAULT 'Bearer',
			expires_at TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return err
}

func CreateTwitchMessagesTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS twitch_messages (
			id SERIAL PRIMARY KEY,
			username VARCHAR(25) NOT NULL,
			display_name VARCHAR(25),
			message TEXT NOT NULL,
			color VARCHAR(7),
			badges JSONB,
			is_mod BOOLEAN DEFAULT FALSE,
			is_vip BOOLEAN DEFAULT FALSE,
			is_broadcaster BOOLEAN DEFAULT FALSE,
			time TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_twitch_messages_time ON twitch_messages(time);
		CREATE INDEX IF NOT EXISTS idx_twitch_messages_username ON twitch_messages(username);
	`)
	return err
}
