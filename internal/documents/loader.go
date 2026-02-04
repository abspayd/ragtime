package documents

import (
	"context"
	"os"
	"path/filepath"
)

func Load(ctx context.Context, path string) ([]byte, error) {

	if filepath.Ext(path) == "pdf" {
		panic("unimplemented")
	}

	return os.ReadFile(path)
}

func listDir(path string) []string {
	if isGitDir(path) {
		return listGitDir(path)
	}

	// TODO: get each file in the path (respect ignore glob if provided)

	return nil
}

func isGitDir(path string) bool {
	// TODO
	return false
}

func listGitDir(path string) []string {
	// TODO
	return nil
}
