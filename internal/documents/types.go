package documents

import "github.com/abspayd/ragtime/internal/documents/splitters"

type Document struct {
	Path   string
	Chunks []splitters.Chunk
}

type ChunkBatch struct {
	Index  int
	Chunks []splitters.Chunk
}

type EmbeddingsBatch struct {
	Index      int
	Embeddings [][]float32
	Chunks     []splitters.Chunk
}

type VectorStoreBatch struct {
	EmbeddingsBatch EmbeddingsBatch
	Path            string
}
