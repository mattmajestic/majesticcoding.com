package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/services"
)

// GetLaLigaSchedule godoc
// @Summary Get La Liga weekly schedule
// @Description Returns matches for the next 7 days from La Liga
// @Tags La Liga
// @Success 200 {object} models.LaLigaScheduleResponse
// @Failure 500 {object} map[string]string
// @Router /laliga/schedule [get]
func GetLaLigaSchedule(c *gin.Context) {
	matches, err := services.FetchLaLigaSchedule()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"matches": matches,
		"count":   len(matches),
	})
}

// LaLigaWidget renders the La Liga widget template
func LaLigaWidget(c *gin.Context) {
	c.HTML(http.StatusOK, "laliga.tmpl", nil)
}