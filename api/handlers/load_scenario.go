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

func LoadScenarios(c *gin.Context) {
	db, err := sql.Open("libsql", os.Getenv("TURSO_DATABASE_URL")+"?authToken="+os.Getenv("TURSO_TOKEN"))
	if err != nil {
		fmt.Println("DB connection error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB connection failed"})
		return
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, user_id, project_name, cloud_provider, created_at FROM scenarios`)
	if err != nil {
		fmt.Println("Query error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query scenarios"})
		return
	}
	defer rows.Close()

	var scenarios []models.Scenario
	for rows.Next() {
		var s models.Scenario
		if err := rows.Scan(&s.ID, &s.UserID, &s.ProjectName, &s.CloudProvider, &s.CreatedAt); err != nil {
			fmt.Println("Scan error:", err)
			continue
		}
		scenarios = append(scenarios, s)
	}

	c.JSON(http.StatusOK, scenarios)
}
