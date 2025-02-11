package dogapi

import "context"

type (
	Client interface {
		GetRandomFunFact(ctx context.Context) (string, error)
	}
)
