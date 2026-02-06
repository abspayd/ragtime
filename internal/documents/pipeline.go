package documents

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/abspayd/ragtime/internal/logger"
	"github.com/abspayd/ragtime/internal/models"
	"github.com/qdrant/go-client/qdrant"
)

func UploadDocuments(paths []string, collection string, embeddingsClient *models.OpenAIClient, qdrantClient *qdrant.Client) error {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if embeddingsClient == nil {
		err := errors.New("Attempted to upload documents without an embeddings model client")
		logger.Log.Error(err.Error())
		return err
	}

	if qdrantClient == nil {
		err := errors.New("Attempted to upload documents without a vector store client")
		logger.Log.Error(err.Error())
		return err
	}

	qdrant_collection_exists, err := qdrantClient.CollectionExists(context.Background(), collection)
	if err != nil {
		return err
	}

	if !qdrant_collection_exists {

		err := qdrantClient.CreateCollection(
			ctx,
			&qdrant.CreateCollection{
				CollectionName: collection,
				VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
					Size:     768, // TODO: validate this against what the embedding model returns
					Distance: qdrant.Distance_Cosine,
				}),
			})
		if err != nil {
			return err
		}
	}

	for _, path := range paths {
		data, err := Load(ctx, path)
		if err != nil {
			return err
		}
		logger.Log.Info("loaded file")

		chunks, err := TextChunks(data)
		if err != nil {
			return err
		}
		logger.Log.Info("extracted text chunks")

		for _, chunk := range chunks {
			embeddings, err := embeddingsClient.Embed(ctx, string(chunk.Text))
			if err != nil {
				return err
			}
			logger.Log.Info("retrieved embeddings")

			for _, embedding := range embeddings.Data {
				payload := map[string]*qdrant.Value{
					"path": {Kind: &qdrant.Value_StringValue{StringValue: path}},
					"text": {Kind: &qdrant.Value_StringValue{StringValue: string(chunk.Text)}},
				}
				response, err := Store(collection, chunk.Index, embedding, payload, qdrantClient)
				if err != nil {
					return err
				}

				logger.Log.Info("stored embeddings", "response", response)
			}
		}
	}

	info, err := qdrantClient.GetCollectionInfo(ctx, collection)
	if err != nil {
		return err
	}
	logger.Log.Info("collection info", "info", info)

	qdrantClient.Close()

	fmt.Println("Done")

	return nil
}
