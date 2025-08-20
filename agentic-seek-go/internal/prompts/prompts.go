package prompts

import (
	"os"
	"path/filepath"
)

func Get(name string) (string, error) {
	// This assumes the executable is run from the `agentic-seek-go` directory.
	path := filepath.Join("prompts", name+".txt")
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
