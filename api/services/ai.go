package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type AIProvider string

const (
	ProviderAnthropic AIProvider = "anthropic"
	ProviderGemini    AIProvider = "gemini"
	ProviderOpenAI    AIProvider = "openai"
	ProviderGroq      AIProvider = "groq"
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

// AnthropicRequest represents the request format for Claude API
type AnthropicRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

// GeminiRequest represents the request format for Gemini API
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

// OpenAIRequest represents the request format for OpenAI API
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
	// Default to Gemini if no provider specified
	if req.Provider == "" {
		req.Provider = ProviderGemini
	}

	// Enhance prompt with RAG context
	enhancedPrompt, err := enhancePromptWithRAG(req.Prompt)
	if err != nil {
		// If RAG fails, log but continue with original prompt
		fmt.Printf("RAG enhancement failed: %v\n", err)
		enhancedPrompt = req.Prompt
	}

	// Update the request with enhanced prompt
	enhancedReq := req
	enhancedReq.Prompt = enhancedPrompt

	switch req.Provider {
	case ProviderAnthropic:
		return callAnthropic(enhancedReq)
	case ProviderGemini:
		return callGemini(enhancedReq)
	case ProviderOpenAI:
		return callOpenAI(enhancedReq)
	case ProviderGroq:
		return callGroq(enhancedReq)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
	}
}

// enhancePromptWithRAG retrieves relevant context and enhances the user prompt
func enhancePromptWithRAG(userPrompt string) (string, error) {
	// Retrieve relevant context (limit to top 3 most relevant pieces)
	contexts, err := RetrieveRelevantContext(userPrompt, 3)
	if err != nil {
		return userPrompt, err
	}

	// Get personality context
	personalityContext := CreatePersonalityContext()

	// Build enhanced prompt
	var enhancedPrompt strings.Builder

	enhancedPrompt.WriteString("You are Matt, a software engineer and content creator. Here's some context about you:\n\n")
	enhancedPrompt.WriteString(personalityContext)
	enhancedPrompt.WriteString("\n\n")

	if len(contexts) > 0 {
		enhancedPrompt.WriteString("Here's some current information about your online presence:\n\n")
		for i, context := range contexts {
			enhancedPrompt.WriteString(fmt.Sprintf("%d. %s\n", i+1, context))
		}
		enhancedPrompt.WriteString("\n")
	}

	enhancedPrompt.WriteString("Please respond as Matt would, using this context to provide helpful and accurate information. ")
	enhancedPrompt.WriteString("If asked about your stats, projects, or online presence, reference the specific numbers and details provided above. ")
	enhancedPrompt.WriteString("Be conversational, helpful, and share insights about software engineering and content creation.\n\n")
	enhancedPrompt.WriteString("User question: ")
	enhancedPrompt.WriteString(userPrompt)

	return enhancedPrompt.String(), nil
}

func callAnthropic(req AIRequest) (*AIResponse, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY not set")
	}

	model := req.Model
	if model == "" {
		model = "claude-3-haiku-20240307" // Cheapest Claude model
	}

	payload := AnthropicRequest{
		Model:     model,
		MaxTokens: 1000,
		Messages: []Message{
			{Role: "user", Content: req.Prompt},
		},
	}

	jsonData, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("anthropic API error: %s", string(body))
	}

	var anthResp AnthropicResponse
	if err := json.Unmarshal(body, &anthResp); err != nil {
		return nil, err
	}

	response := ""
	if len(anthResp.Content) > 0 {
		response = anthResp.Content[0].Text
	}

	return &AIResponse{
		Response: response,
		Provider: string(ProviderAnthropic),
		Model:    model,
	}, nil
}

func callGemini(req AIRequest) (*AIResponse, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}

	model := req.Model
	if model == "" {
		model = "gemini-1.5-flash" // Free/cheap Gemini model
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

func callOpenAI(req AIRequest) (*AIResponse, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not set")
	}

	model := req.Model
	if model == "" {
		model = "gpt-4o-mini" // Cheapest GPT-4 model
	}

	payload := OpenAIRequest{
		Model: model,
		Messages: []OpenAIMessage{
			{Role: "user", Content: req.Prompt},
		},
	}

	jsonData, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openai API error: %s", string(body))
	}

	var openaiResp OpenAIResponse
	if err := json.Unmarshal(body, &openaiResp); err != nil {
		return nil, err
	}

	response := ""
	if len(openaiResp.Choices) > 0 {
		response = openaiResp.Choices[0].Message.Content
	}

	return &AIResponse{
		Response: response,
		Provider: string(ProviderOpenAI),
		Model:    model,
	}, nil
}

func callGroq(req AIRequest) (*AIResponse, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY not set")
	}

	model := req.Model
	if model == "" {
		model = "llama3-8b-8192" // Free Groq model
	}

	payload := OpenAIRequest{ // Groq uses OpenAI-compatible format
		Model: model,
		Messages: []OpenAIMessage{
			{Role: "user", Content: req.Prompt},
		},
	}

	jsonData, _ := json.Marshal(payload)
	httpReq, _ := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("groq API error: %s", string(body))
	}

	var groqResp OpenAIResponse
	if err := json.Unmarshal(body, &groqResp); err != nil {
		return nil, err
	}

	response := ""
	if len(groqResp.Choices) > 0 {
		response = groqResp.Choices[0].Message.Content
	}

	return &AIResponse{
		Response: response,
		Provider: string(ProviderGroq),
		Model:    model,
	}, nil
}

// GetAvailableProviders returns list of providers with available API keys
func GetAvailableProviders() []string {
	var providers []string

	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		providers = append(providers, string(ProviderAnthropic))
	}
	if os.Getenv("GEMINI_API_KEY") != "" {
		providers = append(providers, string(ProviderGemini))
	}
	if os.Getenv("OPENAI_API_KEY") != "" {
		providers = append(providers, string(ProviderOpenAI))
	}
	if os.Getenv("GROQ_API_KEY") != "" {
		providers = append(providers, string(ProviderGroq))
	}

	return providers
}

// GetFallbackProvider returns the first available provider
func GetFallbackProvider() AIProvider {
	providers := GetAvailableProviders()
	if len(providers) == 0 {
		return ""
	}

	// Prefer Gemini for free tier, then Groq, then others
	for _, provider := range []string{"gemini", "groq", "anthropic", "openai"} {
		for _, available := range providers {
			if strings.ToLower(available) == provider {
				return AIProvider(available)
			}
		}
	}

	return AIProvider(providers[0])
}
