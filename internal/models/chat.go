package models

import (
	"context"
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

	return &resp, nil
}
