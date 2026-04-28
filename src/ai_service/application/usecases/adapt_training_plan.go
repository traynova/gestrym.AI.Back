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

type AdaptTrainingPlanUseCase struct {
	repo           repositories.AIRepository
	progressClient clients.ProgressClient
	trainingClient clients.TrainingClient
	aiClient       clients.AIClient
}

func NewAdaptTrainingPlanUseCase(
	repo repositories.AIRepository,
	progressClient clients.ProgressClient,
	trainingClient clients.TrainingClient,
	aiClient clients.AIClient,
) *AdaptTrainingPlanUseCase {
	return &AdaptTrainingPlanUseCase{
		repo:           repo,
		progressClient: progressClient,
		trainingClient: trainingClient,
		aiClient:       aiClient,
	}
}

func (uc *AdaptTrainingPlanUseCase) Execute(ctx context.Context, userID uuid.UUID) (*models.AIRecommendation, error) {
	// 1. Fetch user progress
	progress, err := uc.progressClient.GetUserProgress(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch user progress: %w", err)
	}

	// 2. Prepare AI Prompt
	systemPrompt := `You are Gestrym AI Coach. Analyze the user's performance and adapt their training plan.
	Detect plateaus or rapid progress and adjust volume/intensity.
	Return a JSON with:
	- adaptation_summary (string)
	- adjustments (array of objects with 'exercise', 'change')
	- new_volume (string)
	- deload_needed (boolean)
	- coach_notes (string)
	
	Strictly return ONLY JSON.`

	userPrompt := fmt.Sprintf("Analyze performance: %v", progress)

	// 3. Call AI
	aiResponse, err := uc.aiClient.GenerateResponse(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("ai adaptation failed: %w", err)
	}

	// 4. Send adapted plan to training-service
	var plan interface{}
	json.Unmarshal([]byte(aiResponse), &plan)

	err = uc.trainingClient.CreateTrainingPlan(ctx, plan)
	if err != nil {
		return nil, fmt.Errorf("could not send adapted plan to training service: %w", err)
	}

	// 5. Store AIRecommendation
	inputJSON, _ := json.Marshal(map[string]interface{}{"progress": progress})
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
