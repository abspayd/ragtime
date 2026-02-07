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
		err := errors.New("Unable to upload documents: missing embeddings model client")
		logger.Log.Error(err.Error())
		return err
	}

	if qdrantClient == nil {
		err := errors.New("Unable to upload documents: missing vector store client")
		logger.Log.Error(err.Error())
		return err
	}

	err := buildCollection(ctx, collection, embeddingsClient, qdrantClient)
	if err != nil {
		return fmt.Errorf("Failed to validate or build collection \"s\": %w", err)
	}

	for _, path := range paths {
		data, err := Load(ctx, path)
		if err != nil {
			return fmt.Errorf("Unable to load data from path \"%s\": %w", path, err)
		}
		logger.Log.Info("loaded file")

		chunks, err := TextChunks(data)
		if err != nil {
			return fmt.Errorf("Failed to chunk text: %w", err)
		}
		logger.Log.Info("extracted text chunks")

		for _, chunk := range chunks {
			embeddings, err := embeddingsClient.Embed(ctx, string(chunk.Text))
			if err != nil {
				return fmt.Errorf("Failed to embed text: %w", err)
			}
			logger.Log.Info("retrieved embeddings")

			for _, embedding := range embeddings.Data {
				payload := map[string]*qdrant.Value{
					"path": {Kind: &qdrant.Value_StringValue{StringValue: path}},
					"text": {Kind: &qdrant.Value_StringValue{StringValue: string(chunk.Text)}},
				}
				response, err := Store(collection, chunk.Index, embedding, payload, qdrantClient)
				if err != nil {
					return fmt.Errorf("Failed to store chunks: %w", err)
				}

				logger.Log.Info("stored embeddings", "response", response)
			}
		}
	}

	return nil
}

func buildCollection(ctx context.Context, collection string, embeddingsClient *models.OpenAIClient, qdrantClient *qdrant.Client) error {

	if embeddingsClient == nil {
		err := errors.New("Unable to upload documents: missing embeddings model client")
		logger.Log.Error(err.Error())
		return err
	}

	embeddings_vector_size, err := embeddingsClient.EmbeddingsVectorLength(ctx)
	if err != nil {
		return fmt.Errorf("Unable to determine embeddings vector size: %w", err)
	}

	if qdrantClient == nil {
		err := errors.New("Unable to upload documents: missing vector store client")
		logger.Log.Error(err.Error())
		return err
	}

	qdrant_collection_exists, err := qdrantClient.CollectionExists(context.Background(), collection)
	if err != nil {
		return fmt.Errorf("Unable to determine if collection \"%s\" exists: %w", collection, err)
	}

	if qdrant_collection_exists {
		// validate existing collection
		info, err := qdrantClient.GetCollectionInfo(ctx, collection)
		if err != nil {
			return fmt.Errorf("Unable to get info for collection \"%s\": %w", collection, err)
		}
		collection_vector_size := info.Config.Params.VectorsConfig.GetParams().Size
		if collection_vector_size != uint64(embeddings_vector_size) {
			return fmt.Errorf("Collection \"%s\" has a vector size of %d. Embeddings model returned with incompatible size %d", collection, collection_vector_size, embeddings_vector_size)
		}
	} else {
		// create a new collection with the appropriate size
		err := qdrantClient.CreateCollection(
			ctx,
			&qdrant.CreateCollection{
				CollectionName: collection,
				VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
					Size:     uint64(embeddings_vector_size),
					Distance: qdrant.Distance_Cosine,
				}),
			})
		if err != nil {
			return fmt.Errorf("Unable to create vector store collection: %w", err)
		}
	}

	return nil
}
