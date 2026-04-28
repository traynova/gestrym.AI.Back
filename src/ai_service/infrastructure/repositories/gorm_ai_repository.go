package repositories

import (
	"context"
	"gestrym-ai/src/ai_service/domain/repositories"
	"gestrym-ai/src/common/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type aiRepository struct {
	db *gorm.DB
}

func NewAIRepository(db *gorm.DB) repositories.AIRepository {
	return &aiRepository{db: db}
}

func (r *aiRepository) Save(ctx context.Context, recommendation *models.AIRecommendation) error {
	return r.db.WithContext(ctx).Save(recommendation).Error
}

func (r *aiRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.AIRecommendation, error) {
	var recommendations []models.AIRecommendation
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&recommendations).Error
	return recommendations, err
}
