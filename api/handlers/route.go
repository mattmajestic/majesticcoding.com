package handlers

import (
	"github.com/gin-gonic/gin"
	// "majesticcoding.com/api/middleware" // Temporarily disabled
)

func SetupRoutes(router *gin.Engine) {

	// Serve static files and templates
	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	// Handle favicon.ico requests
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Redirect(301, "https://avatars.githubusercontent.com/u/33904170?v=4")
	})

	// Render routes
	router.GET("/", RenderTemplate("index.tmpl"))
	router.GET("/auth", RenderWithClerk("auth.tmpl"))
	router.GET("/docs", RenderTemplate("docs.tmpl"))
	router.GET("/about", RenderTemplate("about.tmpl"))
	router.GET("/dashboard", RenderTemplate("dashboard.tmpl"))
	router.GET("/infrastructure", RenderTemplate("infrastructure.tmpl"))
	router.GET("/certifications", RenderTemplate("certifications.tmpl"))
	router.GET("/live/", StreamHandler)

	/// Streaming Widgets
	router.GET("/widget/chat", RenderTemplate("chat-widget.tmpl"))
	router.GET("/widget/twitch", RenderTemplate("twitch.tmpl"))
	router.GET("/widget/lavalamp", RenderTemplate("lavalamp.tmpl"))
	router.GET("/widget/globe", GlobeWidgetHandler())
	router.GET("/widget/spotify", RenderSpotify("spotify.tmpl"))
	router.GET("/widget/epl", EPLWidget)
	router.GET("/widget/laliga", LaLigaWidget)

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

	/// Football Leagues
	router.GET("/api/epl/schedule", GetPremierLeagueSchedule)
	router.GET("/api/laliga/schedule", GetLaLigaSchedule)

	/// Geocoding and Globe
	router.GET("/api/geocode", GeocodeHandler())
	router.POST("/api/checkin", PostCheckinHandler())
	router.GET("/api/checkins", GetCheckinsHandler())
	router.GET("/api/checkins/recent", RecentCheckinsHandler())
	router.GET("/api/globe", MapHandler())

	// Spotify
	router.GET("/api/spotify/login", SpotifyLogin)
	router.GET("/api/spotify/callback", SpotifyCallback)
	router.GET("/api/spotify/status", SpotifyStatus)
	router.GET("/api/spotify/current", SpotifyCurrent)

	/// Chat with Websockets
	router.GET("/api/chat", GetMessages)
	router.GET("/api/chat/users", ChatUserCount)
	router.GET("/ws/chat", ChatWebSocket)
	router.GET("/ws/twitch", TwitchMessagesHandler)

	/// App Metrics (Stream)
	router.GET("/api/stream/status", StreamStatusHandler)
	router.GET("/api/metrics", MetricsHandler)

	/// LLM API
	router.POST("/api/llm", PostLLM)
	/// Deploy IAC
	router.GET("/api/deploy/:provider", DeployIACHandler)

	/// Cost Estimation
	router.GET("/api/cost/cloudrun", CloudRunCostHandler)
	router.GET("/api/cost/infracost", InfracostHandler)

	// Dev (uncomment for development)
	// router.POST("/api/chat", PostMessage)
	router.GET("/api/certifications/", CertificateList())
	router.GET("/api/certification/:filename", ServeCertificationPDF())

	// Swagger docs
	DocsHandler(router)
	AboutHandler(router)
	RegisterSwagger(router)

	// Protected API group (temporarily disabled)
	// apiGroup := router.Group("/api")
	// apiGroup.Use(middleware.Auth())
}
