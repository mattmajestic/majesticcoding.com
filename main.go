package main

import (
	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/config"
	"majesticcoding.com/api/handlers"
	"majesticcoding.com/api/middleware"
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
	handlers.StartMessageCleanup()
	services.StartTwitchChatFeed("majesticcodingtwitch")

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	router.SetTrustedProxies(nil)

	handlers.SetupRoutes(router)

	router.Run(":8080")
}
