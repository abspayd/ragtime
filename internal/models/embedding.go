package models

import (
	"context"
	"fmt"
)

func (c *OpenAIClient) Embed(ctx context.Context, input string) (*EmbeddingsResponse, error) {
	req := EmbeddingsRequest{
		Input: input,
		Model: c.Model,
	}

	var resp EmbeddingsResponse
	if err := c.post(ctx, "/v1/embeddings", req, &resp); err != nil {
		return nil, fmt.Errorf("An error occurred when embedding text: %w", err)
	}

	return &resp, nil
}

func (c *OpenAIClient) EmbeddingsVectorLength(ctx context.Context) (int, error) {
	response, err := c.Embed(ctx, "test")
	if err != nil {
		return 0, fmt.Errorf("Unable to find embedding size: %w", err)
	}

	if len(response.Data) > 0 {
		return len(response.Data[0].Embedding), nil
	}

	return 0, fmt.Errorf("Failed to determine size of embeddings vector.")
}
