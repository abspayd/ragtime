package splitters

import (
	"bytes"

	"github.com/abspayd/ragtime/internal/grammars/markdown"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

const MAX_CHUNK_LENGTH = 1024

func TextChunks(text []byte) ([]Chunk, error) {

	chunks := []Chunk{}
	splitText(text, []byte("\n\n"), &chunks)

	return chunks, nil
}

func splitText(text []byte, delimeter []byte, chunks *[]Chunk) {
	text_splits := bytes.SplitAfterSeq(text, delimeter)

	for split := range text_splits {
		if len(split) <= MAX_CHUNK_LENGTH {
			*chunks = append(*chunks, Chunk{
				Index: len(*chunks),
				Text:  string(split),
			})
		} else {
			var next_delim []byte
			switch string(delimeter) {
			case "\n\n":
				next_delim = []byte("\n")
			case "\n":
				next_delim = []byte(" ")
			default:
				next_delim = []byte("")
			}
			splitText(split, next_delim, chunks)
		}
	}
}

func MarkdownChunks(text []byte) ([]Chunk, error) {

	parser := tree_sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(tree_sitter.NewLanguage(markdown.MarkdownLanguage()))

	tree := parser.Parse(text, nil)
	defer tree.Close()

	root := tree.RootNode()

	chunks := make([]Chunk, 0)
	splitTree(root, text, &chunks)

	return chunks, nil
}

// splitTree recursively retrieves chunks from text using a tree sitter parser.
// Any tree node that exceeds MAX_CHUNK_LENGTH will be split into smaller chunks.
//
// If a node exceeds MAX_CHUNK_LENGTH and has no children, plain text splitting will
// be applied through splitText instead.
func splitTree(node *tree_sitter.Node, text []byte, chunks *[]Chunk) {
	if node == nil {
		return
	}

	bytes_start, bytes_end := node.ByteRange()
	if bytes_end-bytes_start < MAX_CHUNK_LENGTH {

		*chunks = append(*chunks, Chunk{
			Index: len(*chunks),
			Text:  string(text[bytes_start:bytes_end]),
		})
		return
	}

	children := node.Children(node.Walk())
	if len(children) == 0 {
		// fallback to plain text splitting if there are no more child nodes
		splitText(text, []byte("\n\n"), chunks)
	}

	for _, child := range node.Children(node.Walk()) {
		splitTree(&child, text, chunks)
	}
}
