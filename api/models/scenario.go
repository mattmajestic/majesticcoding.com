package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Scenario struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	ProjectName   string `json:"project_name"`
	CloudProvider string `json:"cloud_provider"`
	CreatedAt     string `json:"created_at"`
}

// NewScenario initializes a new Scenario instance
func NewScenario(userID, projectName, cloudProvider string) *Scenario {
	return &Scenario{
		ID:            uuid.New().String(),
		UserID:        userID,
		ProjectName:   projectName,
		CloudProvider: cloudProvider,
		CreatedAt:     time.Now().Format(time.RFC3339),
	}
}

// Save inserts the scenario into the DB
func (s *Scenario) Save(db *sql.DB) error {
	stmt := `INSERT INTO scenarios (id, user_id, project_name, cloud_provider, created_at)
	         VALUES (?, ?, ?, ?, ?)`

	_, err := db.Exec(stmt, s.ID, s.UserID, s.ProjectName, s.CloudProvider, s.CreatedAt)
	if err != nil {
		fmt.Println("Failed SQL Exec:", stmt)
		fmt.Printf("Values: %s, %s, %s, %s, %s\n", s.ID, s.UserID, s.ProjectName, s.CloudProvider, s.CreatedAt)
		return err
	}
	return nil
}
