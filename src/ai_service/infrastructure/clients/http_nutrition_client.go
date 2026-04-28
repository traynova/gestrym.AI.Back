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

type nutritionClient struct {
	baseURL string
	client  *http.Client
}

func NewNutritionClient() clients.NutritionClient {
	return &nutritionClient{
		baseURL: os.Getenv("NUTRITION_SERVICE_URL"),
		client:  &http.Client{},
	}
}

func (c *nutritionClient) CreateMealPlan(ctx context.Context, plan interface{}) error {
	url := fmt.Sprintf("%s/nutrition/plans", c.baseURL)
	body, _ := json.Marshal(plan)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("nutrition service returned status %d", resp.StatusCode)
	}

	return nil
}
