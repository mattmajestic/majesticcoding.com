package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type AIProvider string

const (
	ProviderGemini      AIProvider = "gemini"
	ProviderHuggingFace AIProvider = "HuggingFace"
)

type AIRequest struct {
	Prompt   string     `json:"prompt"`
	Provider AIProvider `json:"provider,omitempty"`
	Model    string     `json:"model,omitempty"`
}

type AIResponse struct {
	Response string `json:"response"`
	Provider string `json:"provider"`
	Model    string `json:"model"`
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type OpenAIRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func GenerateAIResponse(req AIRequest) (*AIResponse, error) {
	if req.Provider == "" {
		req.Provider = ProviderGemini
	}

	switch req.Provider {
	case ProviderGemini:
		return callGemini(req)
	case ProviderHuggingFace:
		return callHuggingFace(req)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
	}
}

func callGemini(req AIRequest) (*AIResponse, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}

	model := req.Model
	if model == "" {
		model = "gemini-2.5-flash"
	}

	payload := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: req.Prompt},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(payload)
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, apiKey)
	httpReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini API error: %s", string(body))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, err
	}

	response := ""
	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		response = geminiResp.Candidates[0].Content.Parts[0].Text
	}

	return &AIResponse{
		Response: response,
		Provider: string(ProviderGemini),
		Model:    model,
	}, nil
}

func callHuggingFace(req AIRequest) (*AIResponse, error) {
	apiKey := os.Getenv("HF_TOKEN")
	if apiKey == "" {
		return nil, fmt.Errorf("HF_TOKEN not set")
	}

	model := req.Model
	if model == "" {
		model = "Qwen/Qwen2.5-72B-Instruct"
	}

	payload := OpenAIRequest{
		Model: model,
		Messages: []OpenAIMessage{
			{Role: "user", Content: req.Prompt},
		},
	}

	jsonData, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequest("POST", "https://router.huggingface.co/v1/chat/completions", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("huggingface API error: %s", string(body))
	}

	var hfResp OpenAIResponse
	if err := json.Unmarshal(body, &hfResp); err != nil {
		return nil, err
	}

	response := ""
	if len(hfResp.Choices) > 0 {
		response = hfResp.Choices[0].Message.Content
	}

	return &AIResponse{
		Response: response,
		Provider: string(ProviderHuggingFace),
		Model:    model,
	}, nil
}

// GetAvailableProviders returns list of providers with available API keys
func GetAvailableProviders() []string {
	var providers []string

	if os.Getenv("GEMINI_API_KEY") != "" {
		providers = append(providers, string(ProviderGemini))
	}
	if os.Getenv("HF_TOKEN") != "" {
		providers = append(providers, string(ProviderHuggingFace))
	}

	return providers
}

// GetFallbackProvider returns the first available provider
func GetFallbackProvider() AIProvider {
	providers := GetAvailableProviders()
	if len(providers) == 0 {
		return ""
	}
	return AIProvider(providers[0])
}
