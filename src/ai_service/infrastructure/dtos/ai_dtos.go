package dtos

import "github.com/google/uuid"

type GenerateTrainingPlanRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Level  string    `json:"level" binding:"required"`
}

type GenerateMealPlanRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Goal   string    `json:"goal" binding:"required"`
}

type AdaptPlanRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}
