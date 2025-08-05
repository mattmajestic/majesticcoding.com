package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/app/interfaces"
)

type PricingHandler struct {
	Service interfaces.PricingService
}

func NewPricingHandler(service interfaces.PricingService) *PricingHandler {
	return &PricingHandler{Service: service}
}

func (h *PricingHandler) GetPricing(c *gin.Context) {
	pricing, err := h.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pricing"})
		return
	}
	c.JSON(http.StatusOK, pricing)
}
