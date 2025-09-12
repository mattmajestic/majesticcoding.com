package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/services"
)

// GetPremierLeagueSchedule godoc
// @Summary Get Premier League weekly schedule
// @Description Returns matches for the next 7 days from the Premier League
// @Tags Premier League
// @Success 200 {object} models.PLScheduleResponse
// @Failure 500 {object} map[string]string
// @Router /epl/schedule [get]
func GetPremierLeagueSchedule(c *gin.Context) {
	matches, err := services.FetchPLSchedule()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"matches": matches,
		"count":   len(matches),
	})
}

// EPLWidget renders the EPL widget template
func EPLWidget(c *gin.Context) {
	c.HTML(http.StatusOK, "epl.tmpl", nil)
}