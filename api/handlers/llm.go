package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"majesticcoding.com/api/models"
)

func PostLLM(c *gin.Context) {
	var req models.LLMRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	payload := map[string]string{
		"inputs": req.Prompt, // flan-t5-small only expects 'inputs'
	}
	jsonPayload, _ := json.Marshal(payload)

	hfToken := os.Getenv("HF_TOKEN")
	if hfToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing HF_TOKEN"})
		return
	}

	reqURL := "https://api-inference.huggingface.co/models/mistralai/Mistral-7B-Instruct-v0.3"
	httpReq, _ := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonPayload))
	httpReq.Header.Set("Authorization", "Bearer "+hfToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Request failed"})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":  "Hugging Face API failed",
			"status": resp.StatusCode,
			"body":   string(body),
		})
		return
	}

	var hfResp map[string]interface{}
	if err := json.Unmarshal(body, &hfResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	// Handle both single string and list of results
	text := ""
	if generated, ok := hfResp["generated_text"].(string); ok {
		text = generated
	} else if arr, ok := hfResp["generated_texts"].([]interface{}); ok && len(arr) > 0 {
		text = arr[0].(string)
	}

	c.JSON(http.StatusOK, models.LLMResponse{Response: text})
}
