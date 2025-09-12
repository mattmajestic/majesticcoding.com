package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"majesticcoding.com/api/models"
)

// InsertMessage inserts a new message into the messages table
func InsertMessage(db *sql.DB, content string) error {
	_, err := db.Exec(`INSERT INTO messages (content) VALUES ($1)`, content)
	return err
}

// InsertChatMessage inserts a new chat message with username into the messages table
func InsertChatMessage(db *sql.DB, username, content string) error {
	_, err := db.Exec(`INSERT INTO messages (username, content) VALUES ($1, $2)`, username, content)
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

func InsertCheckin(db *sql.DB, lat, lon float64, city, country string) error {
	_, err := db.Exec(`INSERT INTO checkins (lat, lon, city, country) VALUES ($1, $2, $3, $4)`, lat, lon, city, country)
	return err
}

func GetCheckins(db *sql.DB) ([]models.Checkin, error) {
	rows, err := db.Query(`SELECT id, lat, lon, COALESCE(city, ''), COALESCE(country, ''), checkin_time FROM checkins ORDER BY checkin_time DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checkins []models.Checkin
	for rows.Next() {
		var c models.Checkin
		if err := rows.Scan(&c.ID, &c.Lat, &c.Lon, &c.City, &c.Country, &c.CheckinTime); err != nil {
			return nil, err
		}
		checkins = append(checkins, c)
	}
	return checkins, nil
}

// Spotify token functions
func SaveSpotifyToken(db *sql.DB, accessToken, refreshToken, tokenType string, expiresAt time.Time) error {
	// First, clear any existing tokens (we only store one at a time)
	_, err := db.Exec(`DELETE FROM spotify_tokens`)
	if err != nil {
		return err
	}
	
	// Insert new token
	_, err = db.Exec(`
		INSERT INTO spotify_tokens (access_token, refresh_token, token_type, expires_at) 
		VALUES ($1, $2, $3, $4)
	`, accessToken, refreshToken, tokenType, expiresAt)
	return err
}

func GetSpotifyToken(db *sql.DB) (accessToken, refreshToken, tokenType string, expiresAt time.Time, err error) {
	err = db.QueryRow(`
		SELECT access_token, COALESCE(refresh_token, ''), token_type, expires_at 
		FROM spotify_tokens 
		ORDER BY created_at DESC 
		LIMIT 1
	`).Scan(&accessToken, &refreshToken, &tokenType, &expiresAt)
	return
}

func UpdateSpotifyToken(db *sql.DB, accessToken, refreshToken string, expiresAt time.Time) error {
	_, err := db.Exec(`
		UPDATE spotify_tokens 
		SET access_token = $1, refresh_token = $2, expires_at = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = (SELECT id FROM spotify_tokens ORDER BY created_at DESC LIMIT 1)
	`, accessToken, refreshToken, expiresAt)
	return err
}

func GetRecentCheckins(db *sql.DB, hoursBack int) ([]models.Checkin, error) {
	query := fmt.Sprintf(`
		SELECT id, lat, lon, COALESCE(city, ''), COALESCE(country, ''), checkin_time 
		FROM checkins 
		WHERE checkin_time >= NOW() - INTERVAL '%d hours'
		ORDER BY checkin_time DESC
	`, hoursBack)
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checkins []models.Checkin
	for rows.Next() {
		var c models.Checkin
		if err := rows.Scan(&c.ID, &c.Lat, &c.Lon, &c.City, &c.Country, &c.CheckinTime); err != nil {
			return nil, err
		}
		checkins = append(checkins, c)
	}
	return checkins, nil
}

// InsertTwitchMessage inserts a new Twitch message into the twitch_messages table
func InsertTwitchMessage(db *sql.DB, message models.TwitchMessage) error {
	// Convert badges map to JSON string
	badgesJSON := "{}"
	if len(message.Badges) > 0 {
		badgesBytes, err := json.Marshal(message.Badges)
		if err == nil {
			badgesJSON = string(badgesBytes)
		}
	}
	
	_, err := db.Exec(`
		INSERT INTO twitch_messages (username, display_name, message, color, badges, is_mod, is_vip, is_broadcaster, time) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, message.Username, message.DisplayName, message.Message, message.Color, 
	   badgesJSON, message.IsMod, message.IsVip, message.IsBroadcaster, message.Time)
	return err
}

// GetRecentTwitchMessages fetches the most recent Twitch messages up to the given limit
func GetRecentTwitchMessages(db *sql.DB, limit int) ([]models.TwitchMessage, error) {
	rows, err := db.Query(`
		SELECT id, username, display_name, message, color, badges, is_mod, is_vip, is_broadcaster, time, created_at 
		FROM twitch_messages 
		ORDER BY time DESC 
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.TwitchMessage
	for rows.Next() {
		var msg models.TwitchMessage
		if err := rows.Scan(&msg.ID, &msg.Username, &msg.DisplayName, &msg.Message, &msg.Color, 
						   &msg.Badges, &msg.IsMod, &msg.IsVip, &msg.IsBroadcaster, &msg.Time, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

// Stats storage functions
func InsertYouTubeStats(db *sql.DB, channelID string, subscribers, videos int, views int64) error {
	_, err := db.Exec(`
		INSERT INTO youtube_stats (channel_id, subscriber_count, video_count, view_count) 
		VALUES ($1, $2, $3, $4)
	`, channelID, subscribers, videos, views)
	return err
}

func InsertGitHubStats(db *sql.DB, username string, repos, followers, following, totalStars int) error {
	_, err := db.Exec(`
		INSERT INTO github_stats (username, public_repos, followers, following, total_stars) 
		VALUES ($1, $2, $3, $4, $5)
	`, username, repos, followers, following, totalStars)
	return err
}

func InsertTwitchStats(db *sql.DB, username string, followers, views int, isLive bool) error {
	_, err := db.Exec(`
		INSERT INTO twitch_stats (username, follower_count, view_count, is_live) 
		VALUES ($1, $2, $3, $4)
	`, username, followers, views, isLive)
	return err
}

func InsertLeetCodeStats(db *sql.DB, username string, solved, ranking int, language string) error {
	_, err := db.Exec(`
		INSERT INTO leetcode_stats (username, solved_count, ranking, main_language) 
		VALUES ($1, $2, $3, $4)
	`, username, solved, ranking, language)
	return err
}
