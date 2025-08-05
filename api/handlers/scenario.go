package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"majesticcoding.com/api/models"
)

// SaveScenario stores a scenario
// @Summary Save a scenario
// @Description Save a new scenario to the database
// @Tags Scenarios
// @Accept json
// @Produce json
// @Param scenario body models.Scenario true "Scenario input"
// @Success 200 {object} map[string]string
// @Router /api/scenario [post]

func SaveScenario(c *gin.Context) {
	var input struct {
		UserID        string `json:"user_id"`
		ProjectName   string `json:"project_name"`
		CloudProvider string `json:"cloud_provider"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	scenario := models.NewScenario(input.UserID, input.ProjectName, input.CloudProvider)

	db, err := sql.Open("libsql", os.Getenv("TURSO_DATABASE_URL")+"?authToken="+os.Getenv("TURSO_TOKEN"))
	if err != nil {
		fmt.Println("DB connection error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection failed"})
		return
	}
	defer db.Close()

	if err := scenario.Save(db); err != nil {
		fmt.Println("Insert error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Insert failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Scenario saved", "data": scenario})
}
