package handlers

import (
	"gestrym-ai/src/ai_service/application/usecases"
	"gestrym-ai/src/ai_service/domain/repositories"
	"gestrym-ai/src/ai_service/infrastructure/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type AIHandler struct {
	repo                        repositories.AIRepository
	generateTrainingPlanUseCase *usecases.GenerateTrainingPlanUseCase
	generateMealPlanUseCase     *usecases.GenerateMealPlanUseCase
	adaptTrainingPlanUseCase    *usecases.AdaptTrainingPlanUseCase
	adaptMealPlanUseCase        *usecases.AdaptMealPlanUseCase
}

func NewAIHandler(
	repo repositories.AIRepository,
	generateTrainingPlanUseCase *usecases.GenerateTrainingPlanUseCase,
	generateMealPlanUseCase *usecases.GenerateMealPlanUseCase,
	adaptTrainingPlanUseCase *usecases.AdaptTrainingPlanUseCase,
	adaptMealPlanUseCase *usecases.AdaptMealPlanUseCase,
) *AIHandler {
	return &AIHandler{
		repo:                        repo,
		generateTrainingPlanUseCase: generateTrainingPlanUseCase,
		generateMealPlanUseCase:     generateMealPlanUseCase,
		adaptTrainingPlanUseCase:    adaptTrainingPlanUseCase,
		adaptMealPlanUseCase:        adaptMealPlanUseCase,
	}
}

func (h *AIHandler) GenerateTrainingPlan(c *gin.Context) {
	var req dtos.GenerateTrainingPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.generateTrainingPlanUseCase.Execute(c.Request.Context(), req.UserID, req.Level)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *AIHandler) GenerateMealPlan(c *gin.Context) {
	var req dtos.GenerateMealPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.generateMealPlanUseCase.Execute(c.Request.Context(), req.UserID, req.Goal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *AIHandler) AdaptTrainingPlan(c *gin.Context) {
	var req dtos.AdaptPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.adaptTrainingPlanUseCase.Execute(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AIHandler) AdaptMealPlan(c *gin.Context) {
	var req dtos.AdaptPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.adaptMealPlanUseCase.Execute(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *AIHandler) GetRecommendations(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	res, err := h.repo.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
