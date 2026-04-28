package clients

import (
	"context"
	"fmt"
	"gestrym-ai/src/ai_service/domain/clients"
	"github.com/sashabaranov/go-openai"
	"os"
)

type openAIClient struct {
	client *openai.Client
	model  string
}

func NewOpenAIClient() clients.AIClient {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		// En un entorno real esto debería ser manejado mejor, 
		// pero para inicialización rápida lo dejamos así.
		fmt.Println("WARNING: OPENAI_API_KEY is not set")
	}
	
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = openai.GPT4oMini
	}

	return &openAIClient{
		client: openai.NewClient(apiKey),
		model:  model,
	}
}

func (c *openAIClient) GenerateResponse(ctx context.Context, systemPrompt string, userPrompt string) (string, error) {
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("openai error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}

	return resp.Choices[0].Message.Content, nil
}
