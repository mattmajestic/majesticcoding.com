package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func ServeCertificationPDF() gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")
		filePath := filepath.Join("static", "img", filename)

		c.File(filePath)
	}
}

func CertificateList() gin.HandlerFunc {
	return func(c *gin.Context) {
		filePath := filepath.Join("static", "certifications.json")
		file, err := os.Open(filePath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		defer file.Close()

		var certs interface{}
		if err := json.NewDecoder(file).Decode(&certs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON"})
			return
		}
		c.JSON(http.StatusOK, certs)
	}
}
