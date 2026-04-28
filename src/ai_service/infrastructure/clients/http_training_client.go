package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gestrym-ai/src/ai_service/domain/clients"
	"net/http"
	"os"
)

type trainingClient struct {
	baseURL string
	client  *http.Client
}

func NewTrainingClient() clients.TrainingClient {
	return &trainingClient{
		baseURL: os.Getenv("TRAINING_SERVICE_URL"),
		client:  &http.Client{},
	}
}

func (c *trainingClient) CreateTrainingPlan(ctx context.Context, plan interface{}) error {
	url := fmt.Sprintf("%s/training/plans", c.baseURL)
	body, _ := json.Marshal(plan)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "应用/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("training service returned status %d", resp.StatusCode)
	}

	return nil
}
