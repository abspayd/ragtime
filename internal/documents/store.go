package documents

import (
	"context"
	"os"
	"os/signal"

	"github.com/abspayd/ragtime/internal/logger"
	"github.com/qdrant/go-client/qdrant"
)

func Store(collection string, data []*qdrant.PointStruct, client *qdrant.Client) (*qdrant.UpdateResult, error) {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	logger.Log.Info("Store upserting data")

	result, err := client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collection,
		Points:         data,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
