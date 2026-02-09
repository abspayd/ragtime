package models

import (
	"context"
	"net/http"
)

type Embedder interface {
	Embed(ctx context.Context, input string) ([][]float32, error)
	EmbedBatch(ctx context.Context, input []string) ([][]float32, error)
	VectorSize(ctx context.Context) (int, error)
}

type ChatModel interface {
	Chat(ctx context.Context, messages []string) ([]Message, error)
}

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
	Input string `json:"input"`
	Model string `json:"model"`
}

type EmbeddingsBatchRequest struct {
	Input []string `json:"input"`
	Model string   `json:"model"`
}

type EmbeddingsResponse struct {
	Data  []Embedding `json:"data"`
	Usage Usage       `json:"usage"`
}

type Embedding struct {
	Index     int       `json:"index"`
	Embedding []float32 `json:"embedding"`
}
