package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/models"
	"majesticcoding.com/api/services"
)

func PostLLM(c *gin.Context) {
	var req models.LLMRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// Convert to service request
	aiReq := services.AIRequest{
		Prompt:   req.Prompt,
		Provider: services.AIProvider(req.Provider),
		Model:    req.Model,
	}

	// If no provider specified, use fallback
	if aiReq.Provider == "" {
		aiReq.Provider = services.GetFallbackProvider()
		if aiReq.Provider == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "No AI providers configured. Please set at least one API key: ANTHROPIC_API_KEY, GEMINI_API_KEY, OPENAI_API_KEY, or GROQ_API_KEY",
			})
			return
		}
	}

	// Call AI service
	resp, err := services.GenerateAIResponse(aiReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   "AI service failed",
			"details": err.Error(),
		})
		return
	}

	// Convert to response model
	llmResp := models.LLMResponse{
		Response: resp.Response,
		Provider: resp.Provider,
		Model:    resp.Model,
	}

	c.JSON(http.StatusOK, llmResp)
}

// GetProviders returns available AI providers
func GetProviders(c *gin.Context) {
	providers := services.GetAvailableProviders()
	c.JSON(http.StatusOK, gin.H{
		"providers": providers,
		"fallback":  services.GetFallbackProvider(),
	})
}
