package clients

import (
	"context"
	"github.com/google/uuid"
)

type ProgressClient interface {
	GetUserMetrics(ctx context.Context, userID uuid.UUID) (interface{}, error)
	GetUserProgress(ctx context.Context, userID uuid.UUID) (interface{}, error)
}
