package dogapi

import (
	"context"
	"fmt"
	"net/http"
	"pack-management/internal/pkg/http/client"
)

type (
	dogAIPClient struct {
		baseURL    string
		baseClient client.Client
	}
)

func NewDogAPIClient(baseClient client.Client, baseURL string) Client {
	return &dogAIPClient{
		baseURL:    baseURL,
		baseClient: baseClient,
	}
}

func (c *dogAIPClient) GetRandomFacts(ctx context.Context, limit int) ([]FactResponse, error) {
	if limit < 1 {
		limit = 1
	}

	factList := FactResponseList{}
	err := c.baseClient.Do(
		ctx,
		client.Request{
			Method: http.MethodGet,
			URL:    fmt.Sprintf("%s/facts?limit=%d", c.baseURL, limit),
		},
		&factList,
	)
	if err != nil {
		return nil, err
	}

	return factList.Data, nil
}
