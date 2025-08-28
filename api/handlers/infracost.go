package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/models"
)

// GET /api/cost/infracost?vendor=gcp&service=Cloud%20Run&region=us-central1&purchaseOption=on_demand
func InfracostHandler(c *gin.Context) {
	apiKey := os.Getenv("ICS_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ICS_API_KEY not set"})
		return
	}

	req := models.InfracostRequest{
		Vendor:         c.DefaultQuery("vendor", "gcp"),
		Service:        c.DefaultQuery("service", "Cloud Run"),
		Region:         c.DefaultQuery("region", "us-central1"),
		ProductFamily:  c.Query("productFamily"),
		PurchaseOption: c.DefaultQuery("purchaseOption", "on_demand"),
		Currency:       c.DefaultQuery("currency", "USD"),
	}

	res, err := models.FetchInfracostPrices(apiKey, req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
