package splitters

import "path/filepath"

var SplitterRegistery = map[string]SplitterFunction{
	"":   TextChunks,
	"md": MarkdownChunks,
}

func Register(name string, callback SplitterFunction) {
	if SplitterRegistery[name] != nil {
		return
	}
	SplitterRegistery[name] = callback
}

// GetSplitter retrieves a splitter function for the provided extension name.
// If no splitter was found, the plain text splitter is returned.
func GetSplitter(name string) SplitterFunction {
	value, ok := SplitterRegistery[name]
	if !ok {
		return SplitterRegistery[""]
	}
	return value
}

// SplitterNameForPath returns the splitter name for a path,
// based on its file extension.
func SplitterNameForPath(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 0 {
		return ext[1:]
	}
	return ext
}
