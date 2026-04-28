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

type AdaptMealPlanUseCase struct {
	repo            repositories.AIRepository
	progressClient  clients.ProgressClient
	nutritionClient clients.NutritionClient
	aiClient        clients.AIClient
}

func NewAdaptMealPlanUseCase(
	repo repositories.AIRepository,
	progressClient clients.ProgressClient,
	nutritionClient clients.NutritionClient,
	aiClient clients.AIClient,
) *AdaptMealPlanUseCase {
	return &AdaptMealPlanUseCase{
		repo:            repo,
		progressClient:  progressClient,
		nutritionClient: nutritionClient,
		aiClient:        aiClient,
	}
}

func (uc *AdaptMealPlanUseCase) Execute(ctx context.Context, userID uuid.UUID) (*models.AIRecommendation, error) {
	// 1. Fetch user progress (history)
	progress, err := uc.progressClient.GetUserProgress(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch user progress: %w", err)
	}

	// 2. Prepare AI Prompt for adaptation
	systemPrompt := `You are Gestrym AI Coach. Analyze the user's progress and adapt their meal plan.
	If weight is not decreasing, reduce calories.
	If weight is increasing too fast, adjust macros.
	Return a JSON with:
	- adaptation_summary (string)
	- new_calories (number)
	- new_macros (object)
	- advice (string)
	
	Strictly return ONLY JSON.`

	userPrompt := fmt.Sprintf("Analyze progress: %v", progress)

	// 3. Call AI
	aiResponse, err := uc.aiClient.GenerateResponse(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("ai adaptation failed: %w", err)
	}

	// 4. Send adapted plan to nutrition-service
	var plan interface{}
	json.Unmarshal([]byte(aiResponse), &plan)

	err = uc.nutritionClient.CreateMealPlan(ctx, plan)
	if err != nil {
		return nil, fmt.Errorf("could not send adapted plan to nutrition service: %w", err)
	}

	// 5. Store AIRecommendation
	inputJSON, _ := json.Marshal(map[string]interface{}{"progress": progress})
	recommendation := &models.AIRecommendation{
		UserID: userID,
		Type:   models.NutritionType,
		Input:  string(inputJSON),
		Output: aiResponse,
	}

	err = uc.repo.Save(ctx, recommendation)
	if err != nil {
		return nil, fmt.Errorf("could not save recommendation: %w", err)
	}

	return recommendation, nil
}
