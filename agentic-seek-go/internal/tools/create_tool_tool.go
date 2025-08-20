package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateToolTool is a tool for creating new tools.
type CreateToolTool struct{}

func (t *CreateToolTool) Name() string {
	return "create_tool"
}

func (t *CreateToolTool) Description() string {
	return "Creates a new tool. Args: <file_name> <go_code>"
}

func (t *CreateToolTool) Run(args ...string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("CreateToolTool requires exactly two arguments: the file name and the Go code")
	}
	fileName := args[0]
	code := args[1]

	// In a real application, we would want to be very careful about where we write files.
	// For now, we will write to the tools directory.
	filePath := filepath.Join("agentic-seek-go", "internal", "tools", fileName)

	err := os.WriteFile(filePath, []byte(code), 0644)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Tool '%s' created successfully.", fileName), nil
}

func init() {
	Register(&CreateToolTool{})
}
