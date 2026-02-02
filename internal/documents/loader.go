package documents

import (
	"os"
)

func Load(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	var files []string
	if stat.IsDir() {
		panic("TODO: implement")

		files = listDir(path)
	} else {
		files = []string{path}
	}

	return nil
}

func listDir(path string) []string {

	if isGitDir(path) {
		return listGitDir(path)
	}

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
