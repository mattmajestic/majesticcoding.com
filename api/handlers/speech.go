package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/services"
)

func PostSpeechTranscribe(c *gin.Context) {
	file, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing audio file"})
		return
	}

	uploaded, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to open audio file"})
		return
	}
	defer uploaded.Close()

	text, err := services.TranscribeAudio(uploaded, file.Filename, file.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "transcription failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"text": text})
}
