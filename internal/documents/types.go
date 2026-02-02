package documents

type Document struct {
	Path   string
	Chunks []Chunk
}

type Chunk struct {
	Index int
	Text  []byte
}
