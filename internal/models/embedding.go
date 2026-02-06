package models

import (
	"context"
)

func (c *OpenAIClient) Embed(ctx context.Context, input string) (*EmbeddingsResponse, error) {
	req := EmbeddingsRequest{
		Input: input,
		Model: c.Model,
	}

	var resp EmbeddingsResponse
	if err := c.post(ctx, "/v1/embeddings", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
