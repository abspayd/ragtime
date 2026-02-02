package models

import (
	"context"

	"github.com/abspayd/ragtime/internal/logger"
)

func (c *OpenAIClient) Embed(ctx context.Context, input string) (*EmbeddingsResponse, error) {
	req := EmbeddingsRequest{
		Input:      input,
		Model:      c.Model,
		Dimensions: 3,
	}

	var resp EmbeddingsResponse
	if err := c.post(ctx, "/v1/embeddings", req, &resp); err != nil {
		return nil, err
	}

	logger.Log.Info("Response from server", "response", resp)

	return &resp, nil
}
