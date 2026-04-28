package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"gestrym-ai/src/ai_service/domain/clients"
	"gestrym-ai/src/ai_service/domain/repositories"
	"gestrym-ai/src/common/models"
	"github.com/google/uuid"
)

type GenerateTrainingPlanUseCase struct {
	repo           repositories.AIRepository
	progressClient clients.ProgressClient
	trainingClient clients.TrainingClient
	aiClient       clients.AIClient
}

func NewGenerateTrainingPlanUseCase(
	repo repositories.AIRepository,
	progressClient clients.ProgressClient,
	trainingClient clients.TrainingClient,
	aiClient clients.AIClient,
) *GenerateTrainingPlanUseCase {
	return &GenerateTrainingPlanUseCase{
		repo:           repo,
		progressClient: progressClient,
		trainingClient: trainingClient,
		aiClient:       aiClient,
	}
}

func (uc *GenerateTrainingPlanUseCase) Execute(ctx context.Context, userID uuid.UUID, level string) (*models.AIRecommendation, error) {
	// 1. Fetch user metrics
	metrics, err := uc.progressClient.GetUserMetrics(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch user metrics: %w", err)
	}

	// 2. Prepare AI Prompt
	systemPrompt := `You are Gestrym AI Coach, a world-class personal trainer. 
	Generate a scientific-based training plan in JSON format.
	The JSON must include:
	- frequency (string)
	- focus (string)
	- exercises (array of objects with 'name', 'sets', 'reps', 'rest')
	- coach_notes (string with motivational and technical advice)
	
	Strictly return ONLY JSON.`

	userPrompt := fmt.Sprintf("Generate a plan for a user with these metrics: %v. Level: %s", metrics, level)

	// 3. Call AI
	aiResponse, err := uc.aiClient.GenerateResponse(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("ai generation failed: %w", err)
	}

	// 4. Send plan to training-service
	var plan interface{}
	json.Unmarshal([]byte(aiResponse), &plan)

	err = uc.trainingClient.CreateTrainingPlan(ctx, plan)
	if err != nil {
		return nil, fmt.Errorf("could not send plan to training service: %w", err)
	}

	// 5. Store AIRecommendation
	inputJSON, _ := json.Marshal(map[string]interface{}{"level": level, "metrics": metrics})
	recommendation := &models.AIRecommendation{
		UserID: userID,
		Type:   models.TrainingType,
		Input:  string(inputJSON),
		Output: aiResponse,
	}

	err = uc.repo.Save(ctx, recommendation)
	if err != nil {
		return nil, fmt.Errorf("could not save recommendation: %w", err)
	}

	return recommendation, nil
}
