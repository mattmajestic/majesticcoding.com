package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/models"
)

// GET /api/cost/cloudrun?region=us-central1&currency=USD
func CloudRunCostHandler(c *gin.Context) {
	apiKey := os.Getenv("GCP_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GCP_API_KEY not set"})
		return
	}
	region := c.DefaultQuery("region", "us-central1")
	currency := c.DefaultQuery("currency", "USD")

	out, err := models.FetchCloudRunAverages(apiKey, currency, region)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}
