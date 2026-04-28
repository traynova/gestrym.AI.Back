package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AIRecommendationType string

const (
	TrainingType  AIRecommendationType = "TRAINING"
	NutritionType AIRecommendationType = "NUTRITION"
)

type AIRecommendation struct {
	ID        uuid.UUID            `gorm:"type:uuid;primary_key;" json:"id"`
	UserID    uuid.UUID            `gorm:"type:uuid;not null;index" json:"user_id"`
	Type      AIRecommendationType `gorm:"type:varchar(20);not null" json:"type"`
	Input     string               `gorm:"type:jsonb" json:"input"`  // JSON input context
	Output    string               `gorm:"type:jsonb" json:"output"` // JSON generated plan
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	DeletedAt gorm.DeletedAt       `gorm:"index" json:"-"`
}

func (m *AIRecommendation) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}
