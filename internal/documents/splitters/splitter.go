package splitters

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/abspayd/ragtime/internal/grammars/markdown"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
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

func MarkdownChunks(text []byte) ([]Chunk, error) {

	parser := tree_sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(tree_sitter.NewLanguage(markdown.MarkdownLanguage()))

	tree := parser.Parse(text, nil)
	defer tree.Close()

	root := tree.RootNode()

	fmt.Println(root.ToSexp())

	return nil, nil
}

func ChunkCode(text []byte) ([]Chunk, error) {

	// TODO

	return nil, nil
}
