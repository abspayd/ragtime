package models

import (
	"context"

	"github.com/abspayd/ragtime/internal/logger"
)

func (c *OpenAIClient) Chat(ctx context.Context, messages []Message) (*ChatResponse, error) {
	req := ChatRequest{
		Messages: messages,
		Model:    c.Model,
		Stream:   false,
	}

	var resp ChatResponse
	if err := c.post(ctx, "/v1/chat/completions", req, &resp); err != nil {
		return nil, err
	}

	logger.Log.Info("Response from server", "response", resp)

	return &resp, nil
}
