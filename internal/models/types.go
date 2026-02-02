package models

import "net/http"

type OpenAIClient struct {
	BaseURL    string
	HTTPClient *http.Client
	APIKey     string
	Model      string
}

type ChatRequest struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Stream      bool      `json:"stream,omitempty"`
	Temperature float32   `json:"temperature,omitempty"`
}

type ChatResponse struct {
	Id      string    `json:"id"`
	Model   string    `json:"model"`
	Created int       `json:"created"`
	Usage   Usage     `json:"usage"`
	Choices []Choices `json:"choices"`
}

type Choices struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type EmbeddingsRequest struct {
	Input      string `json:"input"`
	Model      string `json:"model"`
	Dimensions int    `json:"dimensions,omitempty"`
}

type EmbeddingsResponse struct {
	Data []Embedding `json:"data"`
}

type Embedding struct {
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
	Usage     Usage     `json:"usage"`
}

type Data struct {
}
