package qdrant

import (
	"os"

	"github.com/abspayd/ragtime/internal/config"
	"github.com/qdrant/go-client/qdrant"
)

func NewClient() (*qdrant.Client, error) {
	qdrant.NewClient(&qdrant.Config{
		Host:   config.Config.VectorstoreConfig.BaseURL,
		Port:   6664,
		APIKey: os.Getenv("QDRANT_API_KEY"),
	})

	return nil, nil
}
