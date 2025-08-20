package tools

import "fmt"

// Tool is the interface for all tools.
type Tool interface {
	Name() string
	Description() string
	Run(args ...string) (string, error)
}

var registry = make(map[string]Tool)

// Register adds a tool to the registry.
func Register(tool Tool) {
	if _, exists := registry[tool.Name()]; exists {
		// In a real app, you might want to handle this error more gracefully.
		panic(fmt.Sprintf("tool with name '%s' already registered", tool.Name()))
	}
	registry[tool.Name()] = tool
}

// Get retrieves a tool from the registry.
func Get(name string) (Tool, bool) {
	tool, ok := registry[name]
	return tool, ok
}

// List returns a list of all registered tools.
func List() []Tool {
	var tools []Tool
	for _, tool := range registry {
		tools = append(tools, tool)
	}
	return tools
}
