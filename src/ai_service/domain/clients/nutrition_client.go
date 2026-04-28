package clients

import (
	"context"
)

type NutritionClient interface {
	CreateMealPlan(ctx context.Context, plan interface{}) error
}
