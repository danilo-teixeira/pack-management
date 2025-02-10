package pack

import (
	"context"
	"pack-management/internal/pkg/validator"
	"pack-management/internal/services/person"
)

type (
	Service interface {
		CreatePack(ctx context.Context, pack *Entity) (*Entity, error)
	}

	service struct {
		repo          Repository
		personService person.Service
	}

	ServiceParams struct {
		Repo          Repository     `validate:"required"`
		PersonService person.Service `validate:"required"`
	}
)

func NewService(params *ServiceParams) Service {
	params.validate()

	return &service{
		repo:          params.Repo,
		personService: params.PersonService,
	}
}

func (p *ServiceParams) validate() {
	err := validator.ValidateStruct(p)
	if err != nil {
		panic(err)
	}
}

func (s *service) CreatePack(ctx context.Context, pack *Entity) (*Entity, error) {
	var err error
	pack.Sender, err = s.personService.GetOrCreateByName(ctx, pack.Sender.Name)
	if err != nil {
		return nil, err
	}

	pack.Receiver, err = s.personService.GetOrCreateByName(ctx, pack.Receiver.Name)
	if err != nil {
		return nil, err
	}

	pack.Status = StatusCreated

	err = s.repo.Create(ctx, pack)
	if err != nil {
		return nil, err
	}

	go s.setFunFact(ctx, pack)
	go s.setIsHoliday(ctx, pack)

	return pack, nil
}

func (s *service) setFunFact(ctx context.Context, pack *Entity) {
	// TODO: implement
	return
}

func (s *service) setIsHoliday(ctx context.Context, pack *Entity) {
	// TODO: implement
	return
}
