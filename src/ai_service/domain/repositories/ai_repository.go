package repositories

import (
	"context"
	"gestrym-ai/src/common/models"
	"github.com/google/uuid"
)

type AIRepository interface {
	Save(ctx context.Context, recommendation *models.AIRecommendation) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.AIRecommendation, error)
}
