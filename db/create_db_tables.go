package db

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateBronzeSchema(db *sql.DB) error {
	_, err := db.Exec(`CREATE SCHEMA IF NOT EXISTS bronze;`)
	if err != nil {
		return fmt.Errorf("failed to create bronze schema: %w", err)
	}
	return nil
}

func CreateTables(db *sql.DB) error {
	// First ensure bronze schema exists
	if err := CreateBronzeSchema(db); err != nil {
		return err
	}

	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS bronze.example (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL
	);`)
	return err
}

func CreateMessagesTable(db *sql.DB) error {
	// Ensure bronze schema exists
	if err := CreateBronzeSchema(db); err != nil {
		return err
	}

	// First, add username column if it doesn't exist
	_, err := db.Exec(`
		ALTER TABLE bronze.messages ADD COLUMN IF NOT EXISTS username VARCHAR(50)
	`)
	if err != nil {
		log.Printf("Note: Could not add username column (may already exist): %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS bronze.messages (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50),
			content TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return err
}

func CreateCheckinsTable(db *sql.DB) error {
	// Ensure bronze schema exists
	if err := CreateBronzeSchema(db); err != nil {
		return err
	}

	// Create table with correct floating point types (don't drop existing data)
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS bronze.checkins (
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

	return nil
}

func CreateSpotifyTokensTable(db *sql.DB) error {
	// Ensure bronze schema exists
	if err := CreateBronzeSchema(db); err != nil {
		return err
	}

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS bronze.spotify_tokens (
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

func CreateTwitchTokensTable(db *sql.DB) error {
	// Ensure bronze schema exists
	if err := CreateBronzeSchema(db); err != nil {
		return err
	}

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS bronze.twitch_tokens (
			id SERIAL PRIMARY KEY,
			access_token TEXT NOT NULL,
			refresh_token TEXT,
			token_type TEXT DEFAULT 'Bearer',
			expires_at TIMESTAMPTZ NOT NULL,
			scopes TEXT,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return err
}

func CreateTwitchMessagesTable(db *sql.DB) error {
	// Ensure bronze schema exists
	if err := CreateBronzeSchema(db); err != nil {
		return err
	}

	// Drop the unique constraint if it exists
	_, err := db.Exec(`DROP INDEX IF EXISTS bronze.idx_twitch_messages_unique;`)
	if err != nil {
		log.Printf("⚠️  Could not drop unique constraint: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS bronze.twitch_messages (
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
		
		CREATE INDEX IF NOT EXISTS idx_twitch_messages_time ON bronze.twitch_messages(time);
		CREATE INDEX IF NOT EXISTS idx_twitch_messages_username ON bronze.twitch_messages(username);
	`)
	return err
}

func CreateStatsHistoryTables(db *sql.DB) error {
	// Ensure bronze schema exists
	if err := CreateBronzeSchema(db); err != nil {
		return err
	}

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS bronze.youtube_stats (
			id SERIAL PRIMARY KEY,
			channel_id VARCHAR(255),
			subscriber_count INT,
			video_count INT,
			view_count BIGINT,
			recorded_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS bronze.github_stats (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255),
			public_repos INT,
			followers INT,
			following INT,
			total_stars INT,
			recorded_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS bronze.twitch_stats (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255),
			follower_count INT,
			view_count INT,
			is_live BOOLEAN,
			recorded_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS bronze.leetcode_stats (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255),
			solved_count INT,
			ranking INT,
			main_language VARCHAR(100),
			recorded_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return err
}

func CreateTwitchActivitiesTables(db *sql.DB) error {
	// Ensure bronze schema exists
	if err := CreateBronzeSchema(db); err != nil {
		return err
	}

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS bronze.twitch_followers (
			id SERIAL PRIMARY KEY,
			user_id VARCHAR(255) NOT NULL,
			user_login VARCHAR(255) NOT NULL,
			user_name VARCHAR(255) NOT NULL,
			followed_at TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS bronze.twitch_raids (
			id SERIAL PRIMARY KEY,
			from_broadcaster_user_id VARCHAR(255) NOT NULL,
			from_broadcaster_user_login VARCHAR(255) NOT NULL,
			from_broadcaster_user_name VARCHAR(255) NOT NULL,
			to_broadcaster_user_id VARCHAR(255) NOT NULL,
			to_broadcaster_user_login VARCHAR(255) NOT NULL,
			to_broadcaster_user_name VARCHAR(255) NOT NULL,
			viewers INT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS bronze.twitch_subs (
			id SERIAL PRIMARY KEY,
			user_id VARCHAR(255) NOT NULL,
			user_login VARCHAR(255) NOT NULL,
			user_name VARCHAR(255) NOT NULL,
			broadcaster_user_id VARCHAR(255) NOT NULL,
			broadcaster_user_login VARCHAR(255) NOT NULL,
			broadcaster_user_name VARCHAR(255) NOT NULL,
			tier VARCHAR(10) NOT NULL,
			is_gift BOOLEAN DEFAULT FALSE,
			gifter_user_id VARCHAR(255),
			gifter_user_login VARCHAR(255),
			gifter_user_name VARCHAR(255),
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS bronze.twitch_bits (
			id SERIAL PRIMARY KEY,
			user_id VARCHAR(255) NOT NULL,
			user_login VARCHAR(255) NOT NULL,
			user_name VARCHAR(255) NOT NULL,
			broadcaster_user_id VARCHAR(255) NOT NULL,
			broadcaster_user_login VARCHAR(255) NOT NULL,
			broadcaster_user_name VARCHAR(255) NOT NULL,
			is_anonymous BOOLEAN DEFAULT FALSE,
			message TEXT,
			bits INT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_twitch_followers_user_id ON bronze.twitch_followers(user_id);
		CREATE INDEX IF NOT EXISTS idx_twitch_followers_followed_at ON bronze.twitch_followers(followed_at);
		CREATE INDEX IF NOT EXISTS idx_twitch_raids_created_at ON bronze.twitch_raids(created_at);
		CREATE INDEX IF NOT EXISTS idx_twitch_subs_user_id ON bronze.twitch_subs(user_id);
		CREATE INDEX IF NOT EXISTS idx_twitch_subs_created_at ON bronze.twitch_subs(created_at);
		CREATE INDEX IF NOT EXISTS idx_twitch_bits_user_id ON bronze.twitch_bits(user_id);
		CREATE INDEX IF NOT EXISTS idx_twitch_bits_created_at ON bronze.twitch_bits(created_at);
	`)
	return err
}

func CreateUsersTable(db *sql.DB) error {
	// Ensure bronze schema exists
	if err := CreateBronzeSchema(db); err != nil {
		return err
	}

	// First, try to add the raw_data column if it doesn't exist
	_, err := db.Exec(`
		ALTER TABLE bronze.users ADD COLUMN IF NOT EXISTS raw_data JSONB;
	`)
	if err != nil {
		fmt.Printf("Note: Could not add raw_data column (may already exist): %v\n", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS bronze.users (
			id SERIAL PRIMARY KEY,
			supabase_user_id VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) NOT NULL,
			raw_data JSONB NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_users_supabase_id ON bronze.users(supabase_user_id);
		CREATE INDEX IF NOT EXISTS idx_users_email ON bronze.users(email);
		CREATE INDEX IF NOT EXISTS idx_users_raw_data ON bronze.users USING GIN(raw_data);
	`)
	return err
}

func CreateAuthSessionsTable(db *sql.DB) error {
	// Ensure bronze schema exists
	if err := CreateBronzeSchema(db); err != nil {
		return err
	}

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS bronze.auth_sessions (
			token_hash VARCHAR(64) PRIMARY KEY,
			user_data JSONB NOT NULL,
			expires_at TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_auth_sessions_expires ON bronze.auth_sessions(expires_at);
	`)
	return err
}
