package documents

import (
	"context"
	"errors"
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
					Size:     512,
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

			qdrantClient.Upsert(ctx, &qdrant.UpsertPoints{
				CollectionName: collection,
				Wait:           new(bool),
				Points: []*qdrant.PointStruct{
					{
						Id:      qdrant.NewIDNum(uint64(chunk.Index)),
						Payload: qdrant.NewValueMap(map[string]any{"path": path}),
						Vectors: qdrant.NewVectors(embeddings.Data), // TODO: unwrap embedding data and upload to qdrant
					},
				},
				Ordering:         &qdrant.WriteOrdering{},
				ShardKeySelector: &qdrant.ShardKeySelector{},
				UpdateFilter:     &qdrant.Filter{},
			})
			// logger.Log.Info("retrieved embeddings", "embeddings", embeddings.Data)
		}

	}

	return nil
}
