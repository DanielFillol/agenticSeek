package tools

import (
	"fmt"
	"os"
)

// ReadFileTool is a tool for reading files.
type ReadFileTool struct{}

func (t *ReadFileTool) Name() string {
	return "read_file"
}

func (t *ReadFileTool) Description() string {
	return "Reads the content of a file. Args: <file_path>"
}

func (t *ReadFileTool) Run(args ...string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("ReadFileTool requires exactly one argument: the file path")
	}
	content, err := os.ReadFile(args[0])
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// WriteFileTool is a tool for writing to files.
type WriteFileTool struct{}

func (t *WriteFileTool) Name() string {
	return "write_file"
}

func (t *WriteFileTool) Description() string {
	return "Writes content to a file. Args: <file_path> <content>"
}

func (t *WriteFileTool) Run(args ...string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("WriteFileTool requires exactly two arguments: the file path and the content")
	}
	err := os.WriteFile(args[0], []byte(args[1]), 0644)
	if err != nil {
		return "", err
	}
	return "File written successfully.", nil
}

func init() {
	Register(&ReadFileTool{})
	Register(&WriteFileTool{})
}
