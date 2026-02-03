package documents

import (
	"context"
	"os"
)

func Load(ctx context.Context, path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	contents := os.ReadFile(path)

	return nil, nil
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
