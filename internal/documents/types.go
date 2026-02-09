package documents

type Document struct {
	Path   string
	Chunks []Chunk
}

type Chunk struct {
	Index int
	Text  string
}

type EmbeddingsBatch struct {
	Index      int
	Embeddings [][]float32
	Chunks     []Chunk
}

type ChunkBatch struct {
	Index  int
	Chunks []Chunk
}
