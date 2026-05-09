package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if os.Getenv("ENV") == "production" {
		return
	}
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (skipping)")
	}
	// .env.local overrides .env when present
	if err := godotenv.Overload(".env.local"); err == nil {
		log.Println("Loaded .env.local overrides")
	}
}
