package handlers

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func ServeCertificationPDF() gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")
		filePath := filepath.Join("static", "certifications", filename)

		c.File(filePath)
	}
}
