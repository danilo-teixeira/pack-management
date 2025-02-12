package dogapi

import "context"

type (
	Client interface {
		GetRandomFacts(ctx context.Context, limit int) ([]FactResponse, error)
	}

	FactResponseList struct {
		Data []FactResponse `json:"data"`
	}

	FactResponse struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			Body string `json:"body"`
		} `json:"attributes"`
	}
)
