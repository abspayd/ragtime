package markdown

// #cgo CFLAGS: -std=c11 -fPIC
// #include "tree-sitter-markdown/tree-sitter-markdown/src/parser.c"
// #if __has_include("tree-sitter-markdown/tree-sitter-markdown/src/scanner.c")
// #include "tree-sitter-markdown/tree-sitter-markdown/src/scanner.c"
// #endif
import "C"
import "unsafe"

func MarkdownLanguage() unsafe.Pointer {
	return unsafe.Pointer(C.tree_sitter_markdown())
}
