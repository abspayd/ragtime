package documents

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/abspayd/ragtime/internal/logger"
	"github.com/abspayd/ragtime/internal/models"
	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

const maxBatchSize = 32

func UploadDocuments(paths []string, collection string, embedder models.Embedder, qdrantClient *qdrant.Client) error {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if embedder == nil {
		err := errors.New("Unable to upload documents: missing embeddings model client")
		logger.Log.Error(err.Error())
		return err
	}

	if qdrantClient == nil {
		err := errors.New("Unable to upload documents: missing vector store client")
		logger.Log.Error(err.Error())
		return err
	}

	err := buildCollection(ctx, collection, embedder, qdrantClient)
	if err != nil {
		return fmt.Errorf("Failed to validate or build collection \"%s\": %w", collection, err)
	}

	for _, path := range paths {

		// TODO: delete existing data for this path in Qdrant

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

		chunkCh := make(chan ChunkBatch)
		defer close(chunkCh)

		embedCh := make(chan EmbeddingsBatch)
		errCh := make(chan error)

		go embeddingsWorker(ctx, embedder, chunkCh, embedCh, errCh)

		var embeddings_batches []*EmbeddingsBatch

		chunk_batch := ChunkBatch{
			Index:  0,
			Chunks: make([]Chunk, 0),
		}
		for _, chunk := range chunks {
			chunk_batch.Chunks = append(chunk_batch.Chunks, chunk)
			if len(chunk_batch.Chunks) >= maxBatchSize {
				chunkCh <- chunk_batch

				new_embeddings, err := receieveEmbeddings(embedCh, errCh)
				if err != nil {
					return fmt.Errorf("Error embedding chunks: %w", err)
				}

				embeddings_batches = append(embeddings_batches, new_embeddings)

				chunk_batch = ChunkBatch{
					Index:  chunk_batch.Index + 1,
					Chunks: make([]Chunk, 0),
				}
			}
		}

		if len(chunk_batch.Chunks) > 0 {
			chunkCh <- chunk_batch

			new_embeddings, err := receieveEmbeddings(embedCh, errCh)
			if err != nil {
				return fmt.Errorf("Error embedding chunks: %w", err)
			}

			embeddings_batches = append(embeddings_batches, new_embeddings)
		}

		var points []*qdrant.PointStruct
		for _, embedding_batch := range embeddings_batches {

			for index, embedding := range embedding_batch.Embeddings {
				points = append(points, &qdrant.PointStruct{
					Id: qdrant.NewIDUUID(uuid.NewString()),
					Payload: map[string]*qdrant.Value{
						"index": {Kind: &qdrant.Value_IntegerValue{IntegerValue: int64(index)}},
						"path":  {Kind: &qdrant.Value_StringValue{StringValue: path}},
						"text":  {Kind: &qdrant.Value_StringValue{StringValue: string(embedding_batch.Chunks[index].Text)}},
					},
					Vectors: qdrant.NewVectors(embedding...),
				})
			}
		}

		response, err := Store(collection, points, qdrantClient)
		if err != nil {
			return fmt.Errorf("Failed to store chunks: %w", err)
		}

		logger.Log.Debug("UploadDocuments", "response", response)
	}

	return nil
}

func embeddingsWorker(ctx context.Context, embedder models.Embedder, in <-chan ChunkBatch, out chan<- EmbeddingsBatch, errCh chan<- error) {

	for chunk_batch := range in {

		logger.Log.Debug("Embeddings worker receieved chunk", "index", chunk_batch.Index, "size", len(chunk_batch.Chunks))

		chunk_batch_text := make([]string, 0, len(chunk_batch.Chunks))
		for _, text := range chunk_batch.Chunks {
			chunk_batch_text = append(chunk_batch_text, text.Text)
		}

		embeddings, err := embedder.EmbedBatch(ctx, chunk_batch_text)
		if err != nil {
			errCh <- fmt.Errorf("An error occurred during embeddings batch: %w", err)
			continue
		}

		batch_result := EmbeddingsBatch{
			Index:      chunk_batch.Index,
			Embeddings: embeddings,
			Chunks:     chunk_batch.Chunks,
		}

		out <- batch_result
	}

	close(out)
	close(errCh)
}

func receieveEmbeddings(in <-chan EmbeddingsBatch, errCh <-chan error) (*EmbeddingsBatch, error) {
	select {
	case embeddings := <-in:
		return &embeddings, nil
	case err := <-errCh:
		return nil, fmt.Errorf("Error embedding chunks: %w", err)
	}
}

func buildCollection(ctx context.Context, collection string, embedder models.Embedder, qdrantClient *qdrant.Client) error {

	if embedder == nil {
		err := errors.New("Unable to upload documents: missing embeddings model client")
		logger.Log.Error(err.Error())
		return err
	}

	embeddings_vector_size, err := embedder.VectorSize(ctx)
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
