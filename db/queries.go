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

// InsertAIChatMessage inserts a new AI chat exchange into bronze.ai_chat_messages
func InsertAIChatMessage(db *sql.DB, userID, userEmail, provider, model, prompt, response string) error {
	_, err := db.Exec(`
		INSERT INTO ai_chat_messages (supabase_user_id, user_email, provider, model, prompt, response)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, userID, userEmail, provider, model, prompt, response)
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

// CheckCityExists checks if a city already exists in the database
func CheckCityExists(db *sql.DB, city string) (bool, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM checkins WHERE LOWER(city) = LOWER($1)`, city).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
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

// Twitch token functions
func SaveTwitchToken(db *sql.DB, accessToken, refreshToken, tokenType, scopes string, expiresAt time.Time) error {
	// First, clear any existing tokens (we only store one at a time)
	_, err := db.Exec(`DELETE FROM twitch_tokens`)
	if err != nil {
		return err
	}

	// Insert new token
	_, err = db.Exec(`
		INSERT INTO twitch_tokens (access_token, refresh_token, token_type, scopes, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`, accessToken, refreshToken, tokenType, scopes, expiresAt)
	return err
}

func GetTwitchToken(db *sql.DB) (accessToken, refreshToken, tokenType, scopes string, expiresAt time.Time, err error) {
	err = db.QueryRow(`
		SELECT access_token, COALESCE(refresh_token, ''), token_type, COALESCE(scopes, ''), expires_at
		FROM twitch_tokens
		ORDER BY created_at DESC
		LIMIT 1
	`).Scan(&accessToken, &refreshToken, &tokenType, &scopes, &expiresAt)
	return
}

func UpdateTwitchToken(db *sql.DB, accessToken, refreshToken string, expiresAt time.Time) error {
	_, err := db.Exec(`
		UPDATE twitch_tokens
		SET access_token = $1, refresh_token = $2, expires_at = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = (SELECT id FROM twitch_tokens ORDER BY created_at DESC LIMIT 1)
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

func GetLatestLeetCodeStats(db *sql.DB, username string) (*models.LeetCodeStats, error) {
	var stats models.LeetCodeStats
	err := db.QueryRow(`
		SELECT username, solved_count, ranking, main_language
		FROM bronze.leetcode_stats
		WHERE username = $1
		ORDER BY recorded_at DESC
		LIMIT 1
	`, username).Scan(&stats.Username, &stats.SolvedCount, &stats.Ranking, &stats.Languages)

	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// Twitch Activities insert functions
func InsertTwitchFollower(db *sql.DB, follower models.TwitchFollower) error {
	_, err := db.Exec(`
		INSERT INTO twitch_followers (user_id, user_login, user_name, followed_at)
		VALUES ($1, $2, $3, $4)
	`, follower.UserID, follower.UserLogin, follower.UserName, follower.FollowedAt)
	return err
}

func InsertTwitchRaid(db *sql.DB, raid models.TwitchRaid) error {
	_, err := db.Exec(`
		INSERT INTO twitch_raids (from_broadcaster_user_id, from_broadcaster_user_login, from_broadcaster_user_name,
								  to_broadcaster_user_id, to_broadcaster_user_login, to_broadcaster_user_name, viewers)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, raid.FromBroadcasterUserID, raid.FromBroadcasterUserLogin, raid.FromBroadcasterUserName,
		raid.ToBroadcasterUserID, raid.ToBroadcasterUserLogin, raid.ToBroadcasterUserName, raid.Viewers)
	return err
}

func InsertTwitchSub(db *sql.DB, sub models.TwitchSub) error {
	_, err := db.Exec(`
		INSERT INTO twitch_subs (user_id, user_login, user_name, broadcaster_user_id, broadcaster_user_login,
								broadcaster_user_name, tier, is_gift, gifter_user_id, gifter_user_login, gifter_user_name)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, sub.UserID, sub.UserLogin, sub.UserName, sub.BroadcasterUserID, sub.BroadcasterUserLogin,
		sub.BroadcasterUserName, sub.Tier, sub.IsGift, sub.GifterUserID, sub.GifterUserLogin, sub.GifterUserName)
	return err
}

func InsertTwitchBits(db *sql.DB, bits models.TwitchBits) error {
	_, err := db.Exec(`
		INSERT INTO twitch_bits (user_id, user_login, user_name, broadcaster_user_id, broadcaster_user_login,
								broadcaster_user_name, is_anonymous, message, bits)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, bits.UserID, bits.UserLogin, bits.UserName, bits.BroadcasterUserID, bits.BroadcasterUserLogin,
		bits.BroadcasterUserName, bits.IsAnonymous, bits.Message, bits.Bits)
	return err
}

// Get functions for Twitch activities
func GetRecentTwitchFollowers(db *sql.DB, limit int) ([]models.TwitchFollower, error) {
	rows, err := db.Query(`
		SELECT id, user_id, user_login, user_name, followed_at, created_at
		FROM twitch_followers
		ORDER BY followed_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []models.TwitchFollower
	for rows.Next() {
		var f models.TwitchFollower
		if err := rows.Scan(&f.ID, &f.UserID, &f.UserLogin, &f.UserName, &f.FollowedAt, &f.CreatedAt); err != nil {
			return nil, err
		}
		followers = append(followers, f)
	}
	return followers, nil
}

func GetRecentTwitchRaids(db *sql.DB, limit int) ([]models.TwitchRaid, error) {
	rows, err := db.Query(`
		SELECT id, from_broadcaster_user_id, from_broadcaster_user_login, from_broadcaster_user_name,
			   to_broadcaster_user_id, to_broadcaster_user_login, to_broadcaster_user_name, viewers, created_at
		FROM twitch_raids
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var raids []models.TwitchRaid
	for rows.Next() {
		var r models.TwitchRaid
		if err := rows.Scan(&r.ID, &r.FromBroadcasterUserID, &r.FromBroadcasterUserLogin, &r.FromBroadcasterUserName,
			&r.ToBroadcasterUserID, &r.ToBroadcasterUserLogin, &r.ToBroadcasterUserName, &r.Viewers, &r.CreatedAt); err != nil {
			return nil, err
		}
		raids = append(raids, r)
	}
	return raids, nil
}

func GetRecentTwitchSubs(db *sql.DB, limit int) ([]models.TwitchSub, error) {
	rows, err := db.Query(`
		SELECT id, user_id, user_login, user_name, broadcaster_user_id, broadcaster_user_login,
			   broadcaster_user_name, tier, is_gift, gifter_user_id, gifter_user_login, gifter_user_name, created_at
		FROM twitch_subs
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.TwitchSub
	for rows.Next() {
		var s models.TwitchSub
		if err := rows.Scan(&s.ID, &s.UserID, &s.UserLogin, &s.UserName, &s.BroadcasterUserID, &s.BroadcasterUserLogin,
			&s.BroadcasterUserName, &s.Tier, &s.IsGift, &s.GifterUserID, &s.GifterUserLogin, &s.GifterUserName, &s.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, nil
}

func GetRecentTwitchBits(db *sql.DB, limit int) ([]models.TwitchBits, error) {
	rows, err := db.Query(`
		SELECT id, user_id, user_login, user_name, broadcaster_user_id, broadcaster_user_login,
			   broadcaster_user_name, is_anonymous, message, bits, created_at
		FROM twitch_bits
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bitsList []models.TwitchBits
	for rows.Next() {
		var b models.TwitchBits
		if err := rows.Scan(&b.ID, &b.UserID, &b.UserLogin, &b.UserName, &b.BroadcasterUserID, &b.BroadcasterUserLogin,
			&b.BroadcasterUserName, &b.IsAnonymous, &b.Message, &b.Bits, &b.CreatedAt); err != nil {
			return nil, err
		}
		bitsList = append(bitsList, b)
	}
	return bitsList, nil
}

// GetRecentTwitchUsersFromMessages gets unique Twitch users who have sent messages in the past hour
func GetRecentTwitchUsersFromMessages(db *sql.DB, hoursBack int) ([]string, error) {
	query := `
		SELECT DISTINCT username 
		FROM bronze.twitch_messages 
		WHERE time >= NOW() - INTERVAL '%d hours'
		ORDER BY username
	`

	rows, err := db.Query(fmt.Sprintf(query, hoursBack))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		users = append(users, username)
	}
	return users, nil
}
