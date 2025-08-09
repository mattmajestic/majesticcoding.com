package handlers

import (
	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/middleware"
)

func SetupRoutes(router *gin.Engine) {

	// Serve static files and templates
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	// Render routes
	router.GET("/", RenderTemplate("index.tmpl"))
	router.GET("/auth", RenderWithClerk("auth.tmpl"))
	router.GET("/docs", RenderTemplate("docs.tmpl"))
	router.GET("/about", RenderTemplate("about.tmpl"))
	router.GET("/dashboard", RenderTemplate("dashboard.tmpl"))
	router.GET("/certifications", RenderTemplate("certifications.tmpl"))
	router.GET("/live/", StreamHandler)

	// API routes
	/// Scenarios
	router.POST("/api/scenario", SaveScenario)
	router.GET("/api/scenarios", LoadScenarios)

	/// Session Info
	router.GET("/api/session", GetClerkSession)
	router.GET("/api/user/status", AuthStatus)
	router.GET("/user/status", AuthStatusHandler)

	/// 3rd Party APIs (YouTube, Github, Twitch, Leetcode)
	router.GET("/api/stats/:provider", StatsRouter)
	router.GET("/api/git/hash", GitHashHandler)

	/// Chat with Websockets
	router.GET("/api/chat", GetMessages)
	router.GET("/api/chat/users", ChatUserCount)
	router.GET("/ws/chat", ChatWebSocket)

	/// App Metrics (Stream)
	router.GET("/api/stream/status", StreamStatusHandler)
	router.GET("/api/metrics", MetricsHandler)

	/// LLM API
	router.POST("/api/llm", PostLLM)
	/// Deploy IAC
	router.GET("/api/deploy/:provider", DeployIACHandler)

	// Dev (uncomment for development)
	// router.POST("/api/chat", PostMessage)
	router.GET("/api/certification/:filename", ServeCertificationPDF())

	// Swagger docs
	DocsHandler(router)
	RegisterSwagger(router)

	// Protected API group
	apiGroup := router.Group("/api")
	apiGroup.Use(middleware.Auth())
}
