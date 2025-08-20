package llm

import (
	"os"

	"github.com/go-resty/resty/v2"
)

type OpenAI struct {
	client *resty.Client
	model  string
}

type OpenAIChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type OpenAIChatResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

func NewOpenAI(model string) *OpenAI {
	return &OpenAI{
		client: resty.New().SetAuthToken(os.Getenv("OPENAI_API_KEY")),
		model:  model,
	}
}

func (o *OpenAI) Generate(prompt string) (string, error) {
	req := OpenAIChatRequest{
		Model: o.model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	var resp OpenAIChatResponse
	_, err := o.client.R().
		SetBody(req).
		SetResult(&resp).
		Post("https://api.openai.com/v1/chat/completions")

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
