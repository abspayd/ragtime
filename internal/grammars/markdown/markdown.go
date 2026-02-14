package markdown

// #include "tree_sitter/parser.h"
// extern const TSLanguage *tree_sitter_markdown(void);
import "C"
import "unsafe"

func MarkdownLanguage() unsafe.Pointer {
	return unsafe.Pointer(C.tree_sitter_markdown())
}
