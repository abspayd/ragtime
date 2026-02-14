package splitters

type SplitterFunc func(text []byte) ([]Chunk, error)

type Chunk struct {
	Index int
	Text  string
}
