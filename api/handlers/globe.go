package handlers

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/models"
	"majesticcoding.com/api/services"
	"majesticcoding.com/db"
)

func MapHandler() gin.HandlerFunc {
	apiKey := strings.TrimSpace(os.Getenv("GCP_API_KEY"))

	return func(c *gin.Context) {
		if apiKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "missing GCP_API_KEY"})
			return
		}

		q := strings.TrimSpace(c.Query("q"))
		if q == "" {
			q = strings.TrimSpace(c.Query("city"))
		}
		if q == "" {
			q = strings.TrimSpace(c.Query("zip"))
		}

		lat, lng := 0.0, 0.0
		zoom := 3
		label := ""
		title := "Map"

		if v := strings.TrimSpace(c.Query("lat")); v != "" {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				lat = f
			}
		}
		if v := strings.TrimSpace(c.Query("lng")); v != "" {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				lng = f
			}
		}
		if v := strings.TrimSpace(c.Query("zoom")); v != "" {
			if z, err := strconv.Atoi(v); err == nil {
				zoom = z
			}
		}

		// If a text query is provided, prefer geocoding to set center/marker.
		if q != "" {
			if res, err := services.Geocode(c.Request.Context(), q); err == nil && res != nil {
				lat, lng = res.Location.Lat, res.Location.Lng
				label = res.Formatted
				title = res.Formatted
				if zoom == 3 {
					zoom = 12
				}
			}
		}

		data := models.MapPage{
			Title:  title,
			Lat:    lat,
			Lng:    lng,
			Label:  label,
			APIKey: apiKey,
			Zoom:   zoom,
		}
		c.HTML(http.StatusOK, "globe.tmpl", data)
	}
}

func RecentCheckinsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		database := db.GetDB()
		if database == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
			return
		}

		checkins, err := db.GetRecentCheckins(database, 8) // Last 8 hours
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, checkins)
	}
}

func GlobeWidgetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		database := db.GetDB()
		var checkins []models.Checkin
		
		if database != nil {
			if recentCheckins, err := db.GetRecentCheckins(database, 8); err == nil {
				checkins = recentCheckins
			}
		}

		data := gin.H{
			"Title": "Globe Widget",
			"Checkins": checkins,
		}
		c.HTML(http.StatusOK, "globe.tmpl", data)
	}
}
