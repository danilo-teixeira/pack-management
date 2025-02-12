package person

import (
	"context"
	"pack-management/internal/pkg/validator"
)

type (
	Service interface {
		Create(ctx context.Context, person *Entity) error
		GetByName(ctx context.Context, name string) (*Entity, error)
		GetOrCreateByName(ctx context.Context, name string) (*Entity, error)
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

func (s *service) Create(ctx context.Context, person *Entity) error {
	return s.repo.Create(ctx, person)
}

func (s *service) GetByName(ctx context.Context, name string) (*Entity, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *service) GetOrCreateByName(ctx context.Context, name string) (*Entity, error) {
	personEntity, err := s.GetByName(ctx, name)
	if personEntity != nil && err == nil {
		return personEntity, nil
	}

	if err != nil {
		return nil, err
	}

	personEntity = &Entity{
		Name: name,
	}

	err = s.Create(ctx, personEntity)
	if err != nil {
		return nil, err
	}

	return personEntity, nil
}
