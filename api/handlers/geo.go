package handlers

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"strings"

	"majesticcoding.com/api/services"
	"majesticcoding.com/db"

	"github.com/gin-gonic/gin"
)

func GeocodeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		q := strings.TrimSpace(c.Query("q"))
		if q == "" {
			q = strings.TrimSpace(c.Query("city"))
		}
		if q == "" {
			q = strings.TrimSpace(c.Query("zip"))
		}
		if q == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing query; use ?q=, ?city=, or ?zip="})
			return
		}

		// Get username from query parameter (from Twitch chat)
		username := strings.TrimSpace(c.Query("username"))
		if username == "" {
			username = "mattmajestic" // Default fallback
		}

		// Create cache key from query (hash to avoid special characters)
		queryHash := fmt.Sprintf("%x", md5.Sum([]byte(strings.ToLower(q))))
		cacheKey := fmt.Sprintf("geocode:%s", queryHash)

		// Try to get from Redis cache first (24 hours TTL = 86400 seconds)
		cachedJSON, err := services.RedisGetRawJSON(cacheKey)
		if err == nil && cachedJSON != "" {
			log.Printf("✅ Geocode cache HIT for query: %s", q)
			c.Header("Content-Type", "application/json")
			c.String(http.StatusOK, cachedJSON)
			return
		}
		log.Printf("🔍 Geocode cache MISS for query: %s, calling API", q)

		res, err := services.Geocode(c.Request.Context(), q)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		// Cache the geocoding result for 24 hours (86400 seconds)
		if err := services.RedisSetJSON(cacheKey, res, 86400); err != nil {
			log.Printf("⚠️ Failed to cache geocode result: %v", err)
		} else {
			log.Printf("💾 Cached geocode result for 24 hours")
		}

		// Store the result in the database - simple SQL
		database := db.GetDB()
		if database == nil {
			fmt.Printf("ERROR: Database not available\n")
		} else {
			city := res.Components.City
			if city == "" {
				city = res.Formatted
			}
			fmt.Printf("DEBUG: About to check/insert checkin - City: '%s', Country: '%s', Lat: %f, Lng: %f\n",
				city, res.Components.Country, res.Location.Lat, res.Location.Lng)

			// Always insert the checkin (allow duplicates for different users)
			err := db.InsertCheckin(database, res.Location.Lat, res.Location.Lng, city, res.Components.Country)
			if err != nil {
				fmt.Printf("ERROR: Failed to insert checkin: %v\n", err)
			} else {
				fmt.Printf("SUCCESS: New checkin saved for %s from %s!\n", city, username)

				// Get fresh checkins from database and update Redis cache
				checkins, err := db.GetRecentCheckins(database, 8) // Last 8 hours
				if err != nil {
					log.Printf("⚠️ Failed to get recent checkins after insert: %v", err)
					// Fallback: just delete cache to force refresh
					if err := services.RedisDelete("checkins:recent:8h"); err != nil {
						log.Printf("⚠️ Failed to invalidate checkins cache: %v", err)
					}
				} else {
					// Update the cache with fresh data (5 minutes TTL = 300 seconds)
					if err := services.RedisSetJSON("checkins:recent:8h", checkins, 300); err != nil {
						log.Printf("⚠️ Failed to update checkins cache: %v", err)
					} else {
						log.Printf("✅ Updated checkins:recent:8h cache with new %s checkin", city)
					}
				}
			}
		}

		c.JSON(http.StatusOK, res)
	}
}
