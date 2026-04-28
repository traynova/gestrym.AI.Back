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

type GenerateMealPlanUseCase struct {
	repo            repositories.AIRepository
	progressClient  clients.ProgressClient
	nutritionClient clients.NutritionClient
	aiClient        clients.AIClient
}

func NewGenerateMealPlanUseCase(
	repo repositories.AIRepository,
	progressClient clients.ProgressClient,
	nutritionClient clients.NutritionClient,
	aiClient clients.AIClient,
) *GenerateMealPlanUseCase {
	return &GenerateMealPlanUseCase{
		repo:            repo,
		progressClient:  progressClient,
		nutritionClient: nutritionClient,
		aiClient:        aiClient,
	}
}

func (uc *GenerateMealPlanUseCase) Execute(ctx context.Context, userID uuid.UUID, goal string) (*models.AIRecommendation, error) {
	// 1. Fetch user metrics
	metrics, err := uc.progressClient.GetUserMetrics(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch user metrics: %w", err)
	}

	// 2. Prepare AI Prompt
	systemPrompt := `You are Gestrym AI Coach, an expert in high-performance nutrition. 
	Generate a personalized meal plan in JSON format.
	The JSON must include:
	- calories (number)
	- protein (number, grams)
	- carbs (number, grams)
	- fats (number, grams)
	- meals (array of objects with 'name' and 'description')
	- advice (string with a professional coach tip)
	
	Strictly return ONLY JSON.`

	userPrompt := fmt.Sprintf("Generate a plan for a user with these metrics: %v. Goal: %s", metrics, goal)

	// 3. Call AI
	aiResponse, err := uc.aiClient.GenerateResponse(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, fmt.Errorf("ai generation failed: %w", err)
	}

	// 4. Send plan to nutrition-service (assuming it accepts the same JSON structure)
	var plan interface{}
	json.Unmarshal([]byte(aiResponse), &plan)
	
	err = uc.nutritionClient.CreateMealPlan(ctx, plan)
	if err != nil {
		return nil, fmt.Errorf("could not send plan to nutrition service: %w", err)
	}

	// 5. Store AIRecommendation
	inputJSON, _ := json.Marshal(map[string]interface{}{"goal": goal, "metrics": metrics})
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
