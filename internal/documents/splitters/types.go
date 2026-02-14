package splitters

type SplitterFunction func(text []byte) ([]Chunk, error)

type Chunk struct {
	Index int
	Text  string
}
