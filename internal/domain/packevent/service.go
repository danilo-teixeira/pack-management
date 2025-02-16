package packevent

import (
	"context"
	"pack-management/internal/domain/pack"
	"pack-management/internal/pkg/validator"
)

type (
	Service interface {
		CreateEvent(ctx context.Context, event *Entity) (*Entity, error)
	}

	service struct {
		repo        Repository
		packService pack.Service
	}

	ServiceParams struct {
		Repo        Repository   `validate:"required"`
		PackService pack.Service `validate:"required"`
	}
)

func NewService(params *ServiceParams) Service {
	params.validate()

	return &service{
		repo:        params.Repo,
		packService: params.PackService,
	}
}

func (p *ServiceParams) validate() {
	err := validator.ValidateStruct(p)
	if err != nil {
		panic(err)
	}
}

func (s *service) CreateEvent(ctx context.Context, event *Entity) (*Entity, error) {
	_, err := s.packService.GetPackByID(ctx, event.PackID, false)
	if err != nil {
		return nil, err
	}

	err = s.repo.Create(ctx, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}
