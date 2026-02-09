package documents

import (
	"bufio"
	"bytes"
)

func TextChunks(text []byte) ([]Chunk, error) {
	chunks := []Chunk{}

	reader := bytes.NewReader(text)
	scanner := bufio.NewScanner(reader)

	scanner.Split(bufio.ScanLines)

	chunkIndex := 0
	for scanner.Scan() {

		text := scanner.Bytes()

		if len(text) == 0 {
			continue
		}

		chunks = append(chunks, Chunk{
			Index: chunkIndex,
			Text:  scanner.Text(),
		})

		chunkIndex++
	}

	return chunks, nil
}

func MarkdownChunks(path string) {
}

func ChunkCode() {
}
