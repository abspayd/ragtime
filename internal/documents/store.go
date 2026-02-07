package documents

import (
	"context"
	"os"
	"os/signal"

	"github.com/abspayd/ragtime/internal/models"
	"github.com/qdrant/go-client/qdrant"
)

func BulkStore(collection string, embeddings []models.Embedding) {
}

func Store(collection string, id string, embedding models.Embedding, payload map[string]*qdrant.Value, client *qdrant.Client) (*qdrant.UpdateResult, error) {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	result, err := client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collection,
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDUUID(id),
				Payload: payload,
				Vectors: qdrant.NewVectors(embedding.Embedding...),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
