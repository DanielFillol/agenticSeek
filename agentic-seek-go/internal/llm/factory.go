package llm

import "fmt"

func NewLLM(provider, model string) (LLM, error) {
	switch provider {
	case "ollama":
		return NewOllama(model), nil
	case "openai":
		return NewOpenAI(model), nil
	default:
		return nil, fmt.Errorf("unknown LLM provider: %s", provider)
	}
}
