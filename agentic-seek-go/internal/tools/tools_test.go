package tools

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFileTool(t *testing.T) {
	tool := &ReadFileTool{}
	content, err := tool.Run("testdata/read_test.txt")
	if err != nil {
		t.Fatalf("ReadFileTool failed: %v", err)
	}
	expected := "Hello from read_test.txt!"
	if content != expected {
		t.Errorf("ReadFileTool returned wrong content: got '%s', want '%s'", content, expected)
	}
}

func TestWriteFileTool(t *testing.T) {
	tool := &WriteFileTool{}
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "write_test.txt")
	content := "Hello from WriteFileTool!"

	_, err := tool.Run(filePath, content)
	if err != nil {
		t.Fatalf("WriteFileTool failed: %v", err)
	}

	readContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read back test file: %v", err)
	}

	if string(readContent) != content {
		t.Errorf("WriteFileTool wrote wrong content: got '%s', want '%s'", string(readContent), content)
	}
}
