package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/services"
	"majesticcoding.com/db"
)

// RenderSupabaseAuth renders the Supabase authentication page
func RenderSupabaseAuth(c *gin.Context) {
	c.HTML(http.StatusOK, "supabase-auth.tmpl", gin.H{
		"title": "Login - Majestic Coding",
	})
}

// AuthCallbackHandler handles OAuth callbacks
func AuthCallbackHandler(c *gin.Context) {
	// This is where users land after OAuth (Google, GitHub)
	// Supabase handles the token exchange automatically
	c.HTML(http.StatusOK, "supabase-auth.tmpl", gin.H{
		"title":        "Login - Majestic Coding",
		"oauth_return": true,
	})
}

// SupabaseAuthMiddleware verifies Supabase JWT tokens
func SupabaseAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// Verify Supabase JWT token
		user, err := verifySupabaseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Add user info to context
		c.Set("user", user)
		c.Set("user_id", user["sub"])
		c.Set("user_email", user["email"])

		c.Next()
	}
}

// verifySupabaseToken verifies a Supabase JWT token with caching
func verifySupabaseToken(tokenString string) (map[string]interface{}, error) {
	// Clean the token string
	tokenString = strings.TrimSpace(tokenString)

	if tokenString == "" {
		return nil, fmt.Errorf("empty token")
	}

	// Try to get from cache first
	database := db.GetDB()
	if database != nil {
		cachedUser, err := services.GetCachedUserData(database, tokenString)
		if err == nil && cachedUser != nil {
			log.Printf("✅ Auth cache HIT for token")
			return cachedUser, nil
		}
		log.Printf("🔍 Auth cache MISS for token, fetching from Supabase")
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_ANON_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		return nil, fmt.Errorf("Supabase configuration missing")
	}

	// Make request to Supabase to verify token
	url := fmt.Sprintf("%s/auth/v1/user", supabaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+tokenString)
	req.Header.Set("apikey", supabaseKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read the error response for debugging
		var errorBody []byte
		if resp.Body != nil {
			errorBody, _ = json.Marshal(resp.Body)
		}
		log.Printf("Supabase auth failed: status=%d, body=%s", resp.StatusCode, string(errorBody))
		return nil, fmt.Errorf("invalid token: status %d", resp.StatusCode)
	}

	var user map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Printf("Failed to decode Supabase user response: %v", err)
		return nil, fmt.Errorf("failed to decode user data: %v", err)
	}

	// Cache the validated user data
	if database != nil {
		if err := services.SetCachedUserData(database, tokenString, user); err != nil {
			log.Printf("⚠️ Failed to cache user data: %v", err)
		} else {
			log.Printf("💾 Cached user data for token")
		}
	}

	return user, nil
}

// Helper function to get map keys for debugging
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Protected route example
func ProtectedProfileHandler(c *gin.Context) {
	user := c.MustGet("user").(map[string]interface{})
	c.JSON(http.StatusOK, gin.H{
		"message": "This is a protected route!",
		"user":    user,
	})
}

// SupabaseConfigHandler returns Supabase configuration for client
func SupabaseConfigHandler(c *gin.Context) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_ANON_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Supabase configuration missing",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":     supabaseURL,
		"anonKey": supabaseKey,
	})
}

// Get current user info
func GetUserHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer ", "", 1))

	user, err := verifySupabaseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Debug: Print what we got from Supabase
	fmt.Printf("🔍 Raw user data from Supabase: %+v\n", user)

	// Store in database
	database := db.GetDB()
	if database != nil {
		// Check if we have the required fields with nil safety
		var userID, email string
		var ok bool

		// Try different possible field names for user ID
		if userID, ok = user["id"].(string); !ok {
			if userID, ok = user["sub"].(string); !ok {
				if userID, ok = user["user_id"].(string); !ok {
					fmt.Printf("❌ No valid user ID found in token. Available fields: %v\n", getMapKeys(user))
					c.JSON(http.StatusInternalServerError, gin.H{"error": "No user ID in token"})
					return
				}
			}
		}

		if email, ok = user["email"].(string); !ok {
			fmt.Printf("❌ No email found in token. Available fields: %v\n", getMapKeys(user))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No email in token"})
			return
		}

		rawData, _ := json.Marshal(user)
		fmt.Printf("📝 Storing user: ID=%s, Email=%s\n", userID, email)

		_, err = database.Exec(`
			INSERT INTO bronze.users (supabase_user_id, email, raw_data)
			VALUES ($1, $2, $3)
			ON CONFLICT (supabase_user_id) DO UPDATE SET
				email = EXCLUDED.email,
				raw_data = EXCLUDED.raw_data,
				updated_at = CURRENT_TIMESTAMP
		`, userID, email, rawData)

		if err != nil {
			fmt.Printf("❌ DB error: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		} else {
			fmt.Printf("✅ Stored user %s in database\n", email)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// SyncUserHandler syncs Supabase user data to Neon database
func SyncUserHandler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("💥 PANIC in SyncUserHandler: %v\n", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}()

	fmt.Printf("🔧 SyncUserHandler called\n")

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		fmt.Printf("❌ No auth header\n")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token"})
		return
	}

	fmt.Printf("🔑 Auth header present: %s...\n", authHeader[:20])
	tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer ", "", 1))

	fmt.Printf("🔍 Verifying token...\n")
	user, err := verifySupabaseToken(tokenString)
	if err != nil {
		fmt.Printf("❌ Token verification failed: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	fmt.Printf("✅ Token verified for user: %s\n", user["email"])

	database := db.GetDB()
	if database == nil {
		fmt.Printf("❌ Database unavailable\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database unavailable"})
		return
	}

	fmt.Printf("💾 Database connected, syncing...\n")
	err = syncUserToDatabase(database, user)
	if err != nil {
		fmt.Printf("❌ Sync failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Sync failed"})
		return
	}

	fmt.Printf("🎉 Sync completed successfully\n")
	c.JSON(http.StatusOK, gin.H{"message": "User synced"})
}

// SimpleTestHandler displays simple test page
func SimpleTestHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "simple-test.tmpl", gin.H{
		"title": "Simple Test - Majestic Coding",
	})
}

// AutoSyncHandler - visit this URL after login to auto-sync
func AutoSyncHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "auto-sync.tmpl", gin.H{
		"title": "Auto Sync - Majestic Coding",
	})
}

// DebugSyncHandler - simple debug endpoint
func DebugSyncHandler(c *gin.Context) {
	fmt.Printf("🔍 Debug sync called\n")
	fmt.Printf("Auth header: %s\n", c.GetHeader("Authorization"))

	c.JSON(http.StatusOK, gin.H{
		"message": "Debug endpoint works",
		"headers": map[string]string{
			"authorization": c.GetHeader("Authorization"),
			"content-type":  c.GetHeader("Content-Type"),
		},
	})
}

// ShowUserHandler - simple user data display
func ShowUserHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "show-user.tmpl", gin.H{
		"title": "Show User - Majestic Coding",
	})
}

// UserInfoHandler displays user information
func UserInfoHandler(c *gin.Context) {
	// Always render the page - let client-side JavaScript handle auth
	c.HTML(http.StatusOK, "user-info.tmpl", gin.H{})
}

// SettingsHandler displays user settings page
func SettingsHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "signin.tmpl", gin.H{})
}

// SettingsPageHandler displays settings page with clean auth logic
func SettingsPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "settings.tmpl", gin.H{
		"title": "Settings - Majestic Coding",
	})
}

// UserInfoAPIHandler returns user info as JSON (for AJAX calls)
func UserInfoAPIHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	tokenString := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = strings.TrimSpace(authHeader[7:])
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
		return
	}

	user, err := verifySupabaseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Get database connection
	database := db.GetDB()
	var dbUser map[string]interface{}
	if database != nil {
		dbUser, _ = getUserFromDatabase(database, user["sub"].(string))
	}

	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"supabase_user": user,
		"database_user": dbUser,
	})
}

// syncUserToDatabase syncs Supabase user data to Neon database
func syncUserToDatabase(database *sql.DB, user map[string]interface{}) error {
	// Use the same logic as GetUserHandler for consistent field extraction
	var userID, email string
	var ok bool

	// Try different possible field names for user ID
	if userID, ok = user["id"].(string); !ok {
		if userID, ok = user["sub"].(string); !ok {
			if userID, ok = user["user_id"].(string); !ok {
				fmt.Printf("❌ No valid user ID found in sync. Available fields: %v\n", getMapKeys(user))
				return fmt.Errorf("no user ID in token")
			}
		}
	}

	if email, ok = user["email"].(string); !ok {
		fmt.Printf("❌ No email found in sync. Available fields: %v\n", getMapKeys(user))
		return fmt.Errorf("no email in token")
	}

	fmt.Printf("🔄 Syncing user: %s (%s)\n", email, userID)

	// Convert the entire user object to JSON
	rawData, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("❌ JSON marshal error: %v\n", err)
		return err
	}

	// Insert raw JSON data
	query := `
		INSERT INTO bronze.users (supabase_user_id, email, raw_data, created_at, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (supabase_user_id) DO UPDATE SET
			email = EXCLUDED.email,
			raw_data = EXCLUDED.raw_data,
			updated_at = CURRENT_TIMESTAMP
	`

	result, err := database.Exec(query, userID, email, rawData)
	if err != nil {
		fmt.Printf("❌ Database error: %v\n", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("✅ User synced! Rows affected: %d\n", rowsAffected)

	return nil
}

// getUserFromDatabase retrieves user data from Neon database
func getUserFromDatabase(database *sql.DB, userID string) (map[string]interface{}, error) {
	query := `
		SELECT id, supabase_user_id, email, full_name, avatar_url, provider, provider_id,
			   last_sign_in_at, email_verified, phone, user_metadata, app_metadata,
			   created_at, updated_at
		FROM bronze.users
		WHERE supabase_user_id = $1
	`

	row := database.QueryRow(query, userID)

	var (
		id                                               int
		supabaseUserID, email                            string
		fullName, avatarURL, provider, providerID, phone sql.NullString
		lastSignInAt                                     sql.NullTime
		emailVerified                                    bool
		userMetadata, appMetadata                        sql.NullString
		createdAt, updatedAt                             time.Time
	)

	err := row.Scan(&id, &supabaseUserID, &email, &fullName, &avatarURL, &provider, &providerID,
		&lastSignInAt, &emailVerified, &phone, &userMetadata, &appMetadata, &createdAt, &updatedAt)

	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"id":               id,
		"supabase_user_id": supabaseUserID,
		"email":            email,
		"email_verified":   emailVerified,
		"created_at":       createdAt,
		"updated_at":       updatedAt,
	}

	if fullName.Valid {
		result["full_name"] = fullName.String
	}
	if avatarURL.Valid {
		result["avatar_url"] = avatarURL.String
	}
	if provider.Valid {
		result["provider"] = provider.String
	}
	if providerID.Valid {
		result["provider_id"] = providerID.String
	}
	if phone.Valid {
		result["phone"] = phone.String
	}
	if lastSignInAt.Valid {
		result["last_sign_in_at"] = lastSignInAt.Time
	}
	if userMetadata.Valid {
		var metadata map[string]interface{}
		json.Unmarshal([]byte(userMetadata.String), &metadata)
		result["user_metadata"] = metadata
	}
	if appMetadata.Valid {
		var metadata map[string]interface{}
		json.Unmarshal([]byte(appMetadata.String), &metadata)
		result["app_metadata"] = metadata
	}

	return result, nil
}

// Simple Supabase auth status check
func SupabaseAuthStatusHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusOK, gin.H{
			"authenticated": false,
			"message":       "No token provided",
		})
		return
	}

	tokenString := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	user, err := verifySupabaseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"authenticated": false,
			"error":         err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"user_id":       user["sub"],
		"email":         user["email"],
	})
}
