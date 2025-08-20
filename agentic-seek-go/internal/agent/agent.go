package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"flexible-agent/internal/llm"
	"flexible-agent/internal/prompts"
	"flexible-agent/internal/tools"
)

type Agent struct {
	llmClient llm.LLM
}

func NewAgent(llmClient llm.LLM) *Agent {
	return &Agent{
		llmClient: llmClient,
	}
}

func (a *Agent) Plan(request string) ([]string, error) {
	_ = tools.List() // This is to make sure the tools are registered via init()

	availableTools := tools.List()
	var toolDescriptions []string
	for _, tool := range availableTools {
		toolDescriptions = append(toolDescriptions, fmt.Sprintf("- %s: %s", tool.Name(), tool.Description()))
	}

	promptTemplate, err := prompts.Get("plan_prompt")
	if err != nil {
		return nil, fmt.Errorf("failed to get plan prompt: %w", err)
	}

	prompt := fmt.Sprintf(promptTemplate, request, strings.Join(toolDescriptions, "\n"))

	response, err := a.llmClient.Generate(prompt)
	if err != nil {
		return nil, err
	}

	var plan []string
	jsonResponse := extractJSON(response)
	err = json.Unmarshal([]byte(jsonResponse), &plan)
	if err != nil {
		plan = parsePlanFromResponse(response)
		if plan == nil {
			return nil, fmt.Errorf("failed to parse plan from LLM response: %w", err)
		}
	}

	return plan, nil
}

func (a *Agent) Execute(plan []string) error {
	log.Println("Executing plan:")
	for i, step := range plan {
		log.Printf("Step %d: %s\n", i+1, step)

		toolName, args, err := parseToolCall(step)
		if err != nil {
			log.Printf("  Error parsing tool call: %v\n", err)
			continue
		}

		tool, ok := tools.Get(toolName)
		if !ok {
			log.Printf("  Error: tool '%s' not found\n", toolName)
			continue
		}

		result, err := tool.Run(args...)
		if err != nil {
			log.Printf("  Error running tool '%s': %v\n", toolName, err)
			continue
		}

		log.Printf("  Result: %s\n", result)
	}
	return nil
}

func parseToolCall(call string) (string, []string, error) {
	parts := strings.SplitN(call, "(", 2)
	if len(parts) != 2 {
		return "", nil, fmt.Errorf("invalid tool call format: %s", call)
	}
	toolName := parts[0]
	argsStr := strings.TrimRight(parts[1], ")")

	var args []string
	if toolName == "create_tool" {
		// Special handling for create_tool, as the second argument is Go code.
		parts := strings.SplitN(argsStr, ",", 2)
		if len(parts) != 2 {
			return "", nil, fmt.Errorf("invalid create_tool call: expected 2 arguments")
		}
		args = append(args, strings.TrimSpace(strings.Trim(parts[0], `"`)))
		// The go code is the second argument, and it can contain commas.
		args = append(args, parts[1])
	} else {
		if argsStr != "" {
			rawArgs := strings.Split(argsStr, ",")
			for _, arg := range rawArgs {
				args = append(args, strings.TrimSpace(strings.Trim(arg, `"`)))
			}
		}
	}

	return toolName, args, nil
}

// extractJSON extracts the JSON part from a string that might be wrapped in markdown.
func extractJSON(s string) string {
	start := strings.Index(s, "[")
	end := strings.LastIndex(s, "]")
	if start != -1 && end != -1 && start < end {
		return s[start : end+1]
	}
	// also handle ```json ... ```
	start = strings.Index(s, "```json")
	if start != -1 {
		s = s[start+7:]
		end = strings.LastIndex(s, "```")
		if end != -1 {
			s = s[:end]
		}
	}
	// remove newlines
	s = strings.ReplaceAll(s, "\n", "")
	return s
}

func parsePlanFromResponse(response string) []string {
	var plan []string
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "1.") || strings.hasprefix(line, "-") {
			// very basic parsing
			plan = append(plan, strings.trimspace(line[2:]))
		}
	}
	if len(plan) > 0 {
		return plan
	}
	return nil
}
