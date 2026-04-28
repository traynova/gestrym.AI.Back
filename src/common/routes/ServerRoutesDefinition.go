package routes

import (
	"gestrym-ai/docs"
	"gestrym-ai/src/ai_service/application/usecases"
	"gestrym-ai/src/ai_service/infrastructure/clients"
	"gestrym-ai/src/ai_service/infrastructure/handlers"
	"gestrym-ai/src/ai_service/infrastructure/repositories"
	"gestrym-ai/src/common/config"
	"gestrym-ai/src/common/middleware"
	"gestrym-ai/src/common/utils"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type routesDefinition struct {
	serverGroup    *gin.RouterGroup
	publicGroup    *gin.RouterGroup
	privateGroup   *gin.RouterGroup
	internalGroup  *gin.RouterGroup
	protectedGroup *gin.RouterGroup
	logger         utils.ILogger
	aiHandler      *handlers.AIHandler
}

var (
	routesInstance *routesDefinition
	routesOnce     sync.Once
)

func NewRoutesDefinition(serverInstance *gin.Engine) *routesDefinition {
	routesOnce.Do(func() {
		routesInstance = &routesDefinition{}
		routesInstance.logger = utils.NewLogger()
		docs.SwaggerInfo.Title = "Gestrym AI"
		docs.SwaggerInfo.Description = "API para el manejo de entrenamientos."
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.BasePath = "/gestrym-ai"
		routesInstance.addCORSConfig(serverInstance)
		routesInstance.addRoutes(serverInstance)
	})
	return routesInstance
}

func (r *routesDefinition) addCORSConfig(serverInstance *gin.Engine) {
	corsMiddleware := cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
	serverInstance.Use(corsMiddleware)
}

func (r *routesDefinition) addRoutes(serverInstance *gin.Engine) {
	r.addDefaultRoutes(serverInstance)

	// Instantiate DB

	// ── Repositories ─────────────────────────────────────────────────────────
	db := config.NewPostgresConnection().GetDB()
	aiRepo := repositories.NewAIRepository(db)

	// ── Adapters & Services ──────────────────────────────────────────────────
	progressClient := clients.NewProgressClient()
	trainingClient := clients.NewTrainingClient()
	nutritionClient := clients.NewNutritionClient()
	aiClient := clients.NewOpenAIClient()

	// ── Use Cases ────────────────────────────────────────────────────────────
	genTrainingUC := usecases.NewGenerateTrainingPlanUseCase(aiRepo, progressClient, trainingClient, aiClient)
	genMealUC := usecases.NewGenerateMealPlanUseCase(aiRepo, progressClient, nutritionClient, aiClient)
	adaptTrainingUC := usecases.NewAdaptTrainingPlanUseCase(aiRepo, progressClient, trainingClient, aiClient)
	adaptMealUC := usecases.NewAdaptMealPlanUseCase(aiRepo, progressClient, nutritionClient, aiClient)

	// ── Handlers ──────────────────────────────────────────────────────────────
	r.aiHandler = handlers.NewAIHandler(aiRepo, genTrainingUC, genMealUC, adaptTrainingUC, adaptMealUC)

	// ── Router Groups ─────────────────────────────────────────────────────────
	r.serverGroup = serverInstance.Group(docs.SwaggerInfo.BasePath)
	r.serverGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.publicGroup = r.serverGroup.Group("/public")
	r.privateGroup = r.serverGroup.Group("/private")
	r.protectedGroup = r.serverGroup.Group("/protected")

	// Middleware
	r.privateGroup.Use(middleware.SetupJWTMiddleware())
	r.protectedGroup.Use(middleware.SetupApiKeyMiddleware())

	r.addPublicRoutes()
	r.addPrivateRoutes()
	r.addInternalRoutes()
	r.addProtectedRoutes()
}

func (r *routesDefinition) addDefaultRoutes(serverInstance *gin.Engine) {
	serverInstance.GET("/", func(cnx *gin.Context) {
		cnx.JSON(http.StatusOK, map[string]interface{}{
			"code":    "OK",
			"message": "gestrym-nutrition OK...",
			"date":    utils.GetCurrentTime(),
		})
	})

	serverInstance.NoRoute(func(cnx *gin.Context) {
		cnx.JSON(http.StatusNotFound, map[string]interface{}{
			"code":    "NOT_FOUND",
			"message": "Resource not found",
			"date":    utils.GetCurrentTime(),
		})
	})
}

func (r *routesDefinition) addPublicRoutes() {}

func (r *routesDefinition) addPrivateRoutes() {
	ai := r.privateGroup.Group("/ai")
	{
		ai.POST("/generate-training-plan", r.aiHandler.GenerateTrainingPlan)
		ai.POST("/generate-meal-plan", r.aiHandler.GenerateMealPlan)
		ai.POST("/adapt-training-plan", r.aiHandler.AdaptTrainingPlan)
		ai.POST("/adapt-meal-plan", r.aiHandler.AdaptMealPlan)
		ai.GET("/recommendations/:userId", r.aiHandler.GetRecommendations)
	}
}

func (r *routesDefinition) addInternalRoutes() {}

func (r *routesDefinition) addProtectedRoutes() {}
