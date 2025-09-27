package main

import (
	"log"

	"majesticcoding.com/api/config"
	"majesticcoding.com/api/handlers"
	"majesticcoding.com/api/services"
	"majesticcoding.com/db"
)

// @title Majestic Coding API
// @version 1.0
// @description Go API for Full Stack Application
// @host majesticcoding.com
// @BasePath /api

func main() {
	config.LoadEnv()
	handlers.StartBroadcaster()
	db.Connect()

	// Initialize Redis
	if err := services.InitRedis(); err != nil {
		log.Printf("Warning: Failed to initialize Redis: %v", err)
	}

	// Create database tables
	database := db.GetDB()
	if database != nil {
		db.InitializeDatabaseTables(database)
		services.StartSessionCleanup(database)
	}

	handlers.StartMessageCleanup()
	services.StartTwitchChatFeed("majesticcodingtwitch")
	handlers.InitSpotifyClient()

	router := handlers.InitializeRouter()
	router.Run(":8080")
}
