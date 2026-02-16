package splitters

type SplitterFunc func(text []byte, language string) ([]Chunk, error)

type Chunk struct {
	Index int
	Text  string
}
