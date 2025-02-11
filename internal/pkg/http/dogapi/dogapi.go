package dogapi

import "context"


type (
	dogAIPClient struct {
		url string
	}
)

func NewDogAPIClient(url string) Client {
	return &dogAIPClient{
		url: url,
	}
}

func (c *dogAIPClient) GetRandomFunFact(ctx context.Context) (string, error) {
	// TODO: Implement this
	return "Dogs are awesome", nil
}
