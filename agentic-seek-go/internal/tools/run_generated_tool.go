package tools

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// RunGeneratedTool is a tool for running dynamically generated tools.
type RunGeneratedTool struct{}

func (t *RunGeneratedTool) Name() string {
	return "run_generated_tool"
}

func (t *RunGeneratedTool) Description() string {
	return "Runs a dynamically generated tool. Args: <file_name> <tool_args...>"
}

func (t *RunGeneratedTool) Run(args ...string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("RunGeneratedTool requires at least one argument: the file name")
	}
	fileName := args[0]
	toolArgs := args[1:]

	// The path to the tool file is relative to the project root.
	toolPath := filepath.Join("agentic-seek-go", "internal", "tools", fileName)

	// Since I cannot cd into the agentic-seek-go directory, I will assume
	// that the go command is run from the root of the project, and I will
	// use the -C flag to change the directory.
	cmdArgs := []string{"-C", "agentic-seek-go", "run", "internal/tools/" + fileName}
	cmdArgs = append(cmdArgs, toolArgs...)

	cmd := exec.Command("go", cmdArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run generated tool: %s\n%s", err, string(output))
	}

	return string(output), nil
}

func init() {
	Register(&RunGeneratedTool{})
}
