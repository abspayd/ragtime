package models

import (
	"context"
	"fmt"

	"github.com/abspayd/ragtime/internal/logger"
)

func (c *OpenAIClient) Embed(ctx context.Context, input string) ([][]float32, error) {
	req := EmbeddingsRequest{
		Input: input,
		Model: c.Model,
	}

	var resp EmbeddingsResponse
	if err := c.post(ctx, "/v1/embeddings", req, &resp); err != nil {
		return nil, fmt.Errorf("An error occurred when embedding text: %w", err)
	}

	var embeddings [][]float32
	for _, embedding := range resp.Data {
		embeddings = append(embeddings, embedding.Embedding)
	}

	return embeddings, nil
}

func (c *OpenAIClient) EmbedBatch(ctx context.Context, input []string) ([][]float32, error) {
	req := EmbeddingsBatchRequest{
		Input: input,
		Model: c.Model,
	}

	var resp EmbeddingsResponse
	if err := c.post(ctx, "/v1/embeddings", req, &resp); err != nil {
		return nil, fmt.Errorf("An error occurred when embedding text: %w", err)
	}

	var embeddings [][]float32
	for _, embedding := range resp.Data {
		embeddings = append(embeddings, embedding.Embedding)
	}

	return embeddings, nil
}

func (c *OpenAIClient) VectorSize(ctx context.Context) (int, error) {
	response, err := c.Embed(ctx, "test")
	if err != nil {
		return 0, fmt.Errorf("Unable to find embedding size: %w", err)
	}

	logger.Log.Debug("VectorSize", "response", response)

	if len(response) > 0 {
		return len(response[0]), nil
	}

	return 0, fmt.Errorf("Failed to determine size of embeddings vector.")
}
