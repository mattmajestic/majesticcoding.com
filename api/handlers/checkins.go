package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/models"
	"majesticcoding.com/db"
)

// POST /api/checkin
func PostCheckinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var checkin models.Checkin
		if err := c.ShouldBindJSON(&checkin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		database := db.GetDB()
		if database == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database not available"})
			return
		}

		err := db.InsertCheckin(database, checkin.Lat, checkin.Lon, checkin.City, checkin.Country)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save checkin"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// GET /api/checkins
func GetCheckinsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		database := db.GetDB()
		if database == nil {
			c.JSON(http.StatusOK, []models.Checkin{})
			return
		}

		checkins, err := db.GetCheckins(database)
		if err != nil {
			c.JSON(http.StatusOK, []models.Checkin{})
			return
		}

		if checkins == nil {
			checkins = []models.Checkin{}
		}

		c.JSON(http.StatusOK, checkins)
	}
}
