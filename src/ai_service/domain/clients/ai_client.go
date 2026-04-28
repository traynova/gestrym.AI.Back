package clients

import (
	"context"
)

type AIClient interface {
	GenerateResponse(ctx context.Context, systemPrompt string, userPrompt string) (string, error)
}
