package llm

import (
	"github.com/go-resty/resty/v2"
)

type Ollama struct {
	client *resty.Client
	model  string
}

type OllamaChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaChatResponse struct {
	Model     string  `json:"model"`
	CreatedAt string  `json:"created_at"`
	Message   Message `json:"message"`
	Done      bool    `json:"done"`
}

func NewOllama(model string) *Ollama {
	return &Ollama{
		client: resty.New(),
		model:  model,
	}
}

func (o *Ollama) Generate(prompt string) (string, error) {
	req := OllamaChatRequest{
		Model: o.model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	var resp OllamaChatResponse
	_, err := o.client.R().
		SetBody(req).
		SetResult(&resp).
		Post("http://localhost:11434/api/chat")

	if err != nil {
		return "", err
	}

	return resp.Message.Content, nil
}
