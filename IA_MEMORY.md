# IA_MEMORY.md - Gestrym AI Service

## Project Overview
This service is responsible for generating and adapting training and nutrition plans using AI logic. It orchestrates communication between `progress-service`, `training-service`, and `nutrition-service`.

## Architecture
- **Hexagonal Architecture**: Separation of domain, application, and infrastructure.
- **GORM**: Database ORM for PostgreSQL.
- **Dependency Injection**: Used for better testability and decoupling.
- **HTTP Communication**: Communicates with other microservices via REST API.

## Implementation Details

### Models
- **AIRecommendation**: Stores generated plans and their context.
  - `UserID`, `Type` (TRAINING | NUTRITION), `Input` (JSON), `Output` (JSON).

### Use Cases
- `GenerateTrainingPlanUseCase`
- `GenerateMealPlanUseCase`
- `AdaptTrainingPlanUseCase`
- `AdaptMealPlanUseCase`

### Integration Clients
- `TrainingServiceClient`
- `NutritionServiceClient`
- `ProgressServiceClient`

### Endpoints
- `POST /ai/generate-training-plan`
- `POST /ai/generate-meal-plan`
- `POST /ai/adapt-training-plan`
- `POST /ai/adapt-meal-plan`
- `GET /ai/recommendations/:userId`

## Recent Changes
- Initial project structure for `ai-service`.
- Defined `AIRecommendation` model and registered in migrations.
- Implemented `AIRepository` with GORM.
- Implemented HTTP clients for `progress-service`, `training-service`, and `nutrition-service`.
- Implemented use cases for generating and adapting plans.
- Integrated `AIHandler` and registered routes in `ServerRoutesDefinition.go`.
- Verified dependencies with `go mod tidy`.

## Status
- Core AI service functionality: **Completed**
- OpenAI Integration: **Implemented (sashabaranov/go-openai)**
- Integration with other services: **Ready (requires env variables for URLs)**
- Scoring and history: **Implemented**

## OpenAI Integration & "Smart Coach" Algorithm

### 1. Architecture
We use a provider-based architecture with an `AIClient` interface. This allows us to swap OpenAI for any other LLM provider easily.

### 2. Prompt Engineering (The "Brain")
Each use case has a specific **System Prompt** that defines the "Gestrym AI Coach" persona:
- **Nutrition**: Focuses on TDEE, macro distribution, and professional coaching tips.
- **Training**: Focuses on progressive overload, volume management, and technical execution.

### 3. Structured Data Flow
- **Input**: Context harvested from `progress-service` (weight, metrics, history).
- **Processing**: LLM call with `JSON Mode` enabled to ensure the output is machine-readable.
- **Output Orchestration**: The generated JSON is parsed and pushed to `nutrition-service` or `training-service` via HTTP.

### 4. Smart Adaptation Logic
- **Weight Plateau**: The system detects if weight hasn't moved in X weeks and automatically triggers a calorie reduction prompt.
- **Strength Plateau**: If lifting numbers stall, the AI suggests volume adjustments or deload weeks.
