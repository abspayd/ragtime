package documents

import (
	"bufio"
	"os"
)

func TextChunks(path string) ([]Chunk, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	chunkIndex := 0
	chunks := []Chunk{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		chunks = append(chunks, Chunk{
			Index: chunkIndex,
			Text:  scanner.Bytes(),
		})

		chunkIndex++
	}

	return chunks, nil
}

func MarkdownChunks(path string) {
}

func ChunkCode() {
}
