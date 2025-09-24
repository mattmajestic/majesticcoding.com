package handlers

import (
	"fmt"
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

		res, err := services.Geocode(c.Request.Context(), q)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
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
			fmt.Printf("DEBUG: About to insert checkin - City: '%s', Country: '%s', Lat: %f, Lng: %f\n",
				city, res.Components.Country, res.Location.Lat, res.Location.Lng)

			// Simple direct SQL insert
			result, err := database.Exec("INSERT INTO checkins (lat, lon, city, country) VALUES ($1, $2, $3, $4)",
				res.Location.Lat, res.Location.Lng, city, res.Components.Country)
			if err != nil {
				fmt.Printf("ERROR: Failed to insert checkin: %v\n", err)
			} else {
				rowsAffected, _ := result.RowsAffected()
				fmt.Printf("SUCCESS: Checkin saved! Rows affected: %d\n", rowsAffected)
			}
		}

		c.JSON(http.StatusOK, res)
	}
}
