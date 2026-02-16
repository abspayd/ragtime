package splitters

import (
	"bytes"
	"fmt"
	"path/filepath"
	"unsafe"

	"github.com/abspayd/ragtime/internal/documents/grammars/markdown"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
	tree_sitter_php "github.com/tree-sitter/tree-sitter-php/bindings/go"
)

const MAX_CHUNK_LENGTH = 1024

var language_path_extensions = map[string]string{
	".php": "php",
	".go":  "go",
	".md":  "markdown",
	".mdx": "markdown",
	".c":   "c",
	".h":   "c",
}

func TextChunks(text []byte) []Chunk {

	chunks := []Chunk{}
	splitText(text, []byte("\n\n"), &chunks)

	return chunks
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

func Split(text []byte, path string) []Chunk {

	var grammar unsafe.Pointer
	switch language_path_extensions[filepath.Ext(path)] {
	case "go":
		grammar = tree_sitter_go.Language()
	case "markdown":
		fmt.Println("Using markdown splitter")
		grammar = markdown.Language()
	case "php":
		grammar = tree_sitter_php.LanguagePHP()
	}

	if grammar == nil {
		fmt.Println("Using plaintext splitter")
		return TextChunks(text)
	}

	ts_language := tree_sitter.NewLanguage(grammar)

	parser := tree_sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(ts_language)

	tree := parser.Parse(text, nil)
	defer tree.Close()

	root := tree.RootNode()

	chunks := make([]Chunk, 0)
	splitTree(root, text, &chunks)

	return chunks
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
