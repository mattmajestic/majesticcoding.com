package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var Database *sql.DB

func Connect() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("Error: DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error: Database ping failed: %v", err)
	}

	// Set default schema to bronze
	_, err = db.Exec("SET search_path = bronze, public")
	if err != nil {
		log.Printf("Warning: Could not set schema to bronze: %v", err)
	} else {
	}

	Database = db
}

func GetDB() *sql.DB {
	return Database
}
