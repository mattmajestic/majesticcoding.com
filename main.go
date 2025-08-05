package main

import (
	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/config"
	"majesticcoding.com/api/handlers"
	"majesticcoding.com/db"
)

// @title Majestic Coding API
// @version 1.0
// @description Go API for Full Stack Application
// @host localhost:8080
// @BasePath /api

func main() {
	config.LoadEnv()
	handlers.StartBroadcaster()
	db.Connect()

	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1", "majesticcoding.com"})

	handlers.SetupRoutes(router)

	router.Run(":8080")
}
