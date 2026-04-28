package clients

import (
	"context"
)

type TrainingClient interface {
	CreateTrainingPlan(ctx context.Context, plan interface{}) error
}
