package documents

type Document struct {
	Path   string
	Chunks []Chunk
}

type Chunk struct {
	Index int
	Text  string
}

type ChunkBatch struct {
	Index  int
	Chunks []Chunk
}

type EmbeddingsBatch struct {
	Index      int
	Embeddings [][]float32
	Chunks     []Chunk
}

type VectorStoreBatch struct {
	EmbeddingsBatch EmbeddingsBatch
	Path            string
}
