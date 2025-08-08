package models

type LLMRequest struct {
	Prompt  string `json:"prompt"`
	Context string `json:"context"`
}

type LLMResponse struct {
	Response string `json:"response"`
}
