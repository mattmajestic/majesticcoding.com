package handlers

import (
	"github.com/gin-gonic/gin"
	"majesticcoding.com/app/models"
)

func SaveDiagram(c *gin.Context) {
	var diagram models.Diagram
	if err := c.ShouldBindJSON(&diagram); err != nil {
		c.JSON(400, gin.H{"error": "invalid diagram"})
		return
	}

	// TODO: Save to DB or file
	c.JSON(200, gin.H{"status": "saved", "node_count": len(diagram.Nodes)})
}
