package handlers

import (
	"github.com/gin-gonic/gin"
	"majesticcoding.com/app/middleware"
)

func SetupRoutes(router *gin.Engine) {

	// Serve static files and templates
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	// Websockets
	router.GET("/ws/chat", ChatWebSocket)

	// Public routes
	router.GET("/", RenderTemplate("index.tmpl"))
	router.GET("/auth", RenderWithClerk("auth.tmpl"))
	router.GET("/user/status", AuthStatusHandler)
	router.GET("/docs", RenderTemplate("docs.tmpl"))
	router.GET("/dashboard", RenderTemplate("dashboard.tmpl"))
	router.GET("/live/", StreamHandler)
	router.GET("/api/stream/status", StreamStatusHandler)
	router.GET("/api/metrics", MetricsHandler)

	// API routes
	router.POST("/api/scenario", SaveScenario)
	router.GET("/api/scenarios", LoadScenarios)
	router.GET("/api/stats/:provider", StatsRouter)
	router.GET("/api/chat", GetMessages)
	// router.POST("/api/chat", PostMessage)
	// router.GET("/api/certification/:filename", ServeCertificationPDF())

	// Swagger docs
	DocsHandler(router)
	RegisterSwagger(router)

	// Protected API group
	apiGroup := router.Group("/api")
	apiGroup.Use(middleware.Auth())
	// apiGroup.GET("/dashboard", RenderTemplate("dashboard.tmpl"))
}
