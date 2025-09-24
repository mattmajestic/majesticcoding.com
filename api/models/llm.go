package models

type LLMRequest struct {
	Prompt   string `json:"prompt"`
	Context  string `json:"context,omitempty"`
	Provider string `json:"provider,omitempty"`
	Model    string `json:"model,omitempty"`
}

type LLMResponse struct {
	Response string `json:"response"`
	Provider string `json:"provider"`
	Model    string `json:"model"`
}
