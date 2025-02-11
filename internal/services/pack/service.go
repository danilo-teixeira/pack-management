package pack

import (
	"context"
	"fmt"
	"pack-management/internal/pkg/http/dogapi"
	naegerdateapi "pack-management/internal/pkg/http/nagerdateapi"
	"pack-management/internal/pkg/validator"
	"pack-management/internal/services/person"
	"slices"
)

type (
	Service interface {
		CreatePack(ctx context.Context, pack *Entity) (*Entity, error)
	}

	service struct {
		repo               Repository
		personService      person.Service
		dogAPIClient       dogapi.Client
		nagerDateAPIClient naegerdateapi.Client
	}

	ServiceParams struct {
		Repo               Repository           `validate:"required"`
		PersonService      person.Service       `validate:"required"`
		DogAPIClient       dogapi.Client        `validate:"required"`
		NagerDateAPIClient naegerdateapi.Client `validate:"required"`
	}
)

func NewService(params *ServiceParams) Service {
	params.validate()

	return &service{
		repo:               params.Repo,
		personService:      params.PersonService,
		dogAPIClient:       params.DogAPIClient,
		nagerDateAPIClient: params.NagerDateAPIClient,
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
	funFact, err := s.dogAPIClient.GetRandomFunFact(ctx)
	if err != nil {
		// TODO: implement error handler
		fmt.Println(err)
		return
	}

	pack.FunFact = &funFact

	err = s.repo.UpdateByID(ctx, pack.ID, pack)
	if err != nil {
		// TODO: implement error handler
		fmt.Println(err)
	}
}

func (s *service) setIsHoliday(ctx context.Context, pack *Entity) {
	year := pack.EstimatedDeliveryDate[:4]

	holidayResp, err := s.nagerDateAPIClient.GetHolidays(ctx, "BR", year) // TODO: Get from cache or from API???
	if err != nil {
		// TODO: implement error handler
		fmt.Println(err)
	}

	isHoliday := slices.ContainsFunc(holidayResp, func(holiday naegerdateapi.Holiday) bool {
		return holiday.Date == pack.EstimatedDeliveryDate
	})

	pack.IsHoliday = &isHoliday

	err = s.repo.UpdateByID(ctx, pack.ID, pack)
	if err != nil {
		// TODO: implement error handler
		fmt.Println(err)
	}
}
