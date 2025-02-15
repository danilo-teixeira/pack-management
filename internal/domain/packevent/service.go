package packevent

import (
	"context"
	"pack-management/internal/pkg/validator"
)

type (
	Service interface {
		CreateEvent(ctx context.Context, event *Entity) (*Entity, error)
	}

	service struct {
		repo Repository
	}

	ServiceParams struct {
		Repo Repository `validate:"required"`
	}
)

func NewService(params *ServiceParams) Service {
	params.validate()

	return &service{
		repo: params.Repo,
	}
}

func (p *ServiceParams) validate() {
	err := validator.ValidateStruct(p)
	if err != nil {
		panic(err)
	}
}

func (s *service) CreateEvent(ctx context.Context, event *Entity) (*Entity, error) {
	err := s.repo.Create(ctx, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}
