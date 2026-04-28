package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"gestrym-ai/src/ai_service/domain/clients"
	"github.com/google/uuid"
	"net/http"
	"os"
)

type progressClient struct {
	baseURL string
	client  *http.Client
}

func NewProgressClient() clients.ProgressClient {
	return &progressClient{
		baseURL: os.Getenv("PROGRESS_SERVICE_URL"),
		client:  &http.Client{},
	}
}

func (c *progressClient) GetUserMetrics(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	url := fmt.Sprintf("%s/progress/metrics/%s", c.baseURL, userID)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("progress service returned status %d", resp.StatusCode)
	}

	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c *progressClient) GetUserProgress(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	url := fmt.Sprintf("%s/progress/history/%s", c.baseURL, userID)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("progress service returned status %d", resp.StatusCode)
	}

	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}
