package parse

import (
	"os"
	"path/filepath"
)

// Reads a markdown file from the given path and returns its raw bytes and filename.
func GetContent(path string) ([]byte, string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}

	return content, filepath.Base(path), nil
}
