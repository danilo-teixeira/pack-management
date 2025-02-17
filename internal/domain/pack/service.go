package pack

import (
	"context"
	"log"
	"pack-management/internal/domain/holiday"
	"pack-management/internal/domain/person"
	"pack-management/internal/pkg/http/dogapi"
	"pack-management/internal/pkg/pagination"
	"pack-management/internal/pkg/validator"
	"time"
)

type (
	Service interface {
		CreatePack(ctx context.Context, pack *Entity) (*Entity, error)
		ListPacks(ctx context.Context, filters *ListFilters) ([]*Entity, *pagination.Metadata, error)
		GetPackByID(ctx context.Context, id string, withEvents bool) (*Entity, error)
		UpdatePackStatusByID(ctx context.Context, id string, pack *Entity) (*Entity, error)
		CancelPackStatusByID(ctx context.Context, id string) (*Entity, error)
	}

	ListFilters struct {
		SenderName   *string
		ReceiverName *string
		PageSize     int
		PageCursor   *string
	}

	service struct {
		repo           Repository
		personService  person.Service
		dogAPIClient   dogapi.Client
		holidayService holiday.Service
	}

	ServiceParams struct {
		Repo           Repository      `validate:"required"`
		PersonService  person.Service  `validate:"required"`
		DogAPIClient   dogapi.Client   `validate:"required"`
		HolidayService holiday.Service `validate:"required"`
	}
)

func NewService(params *ServiceParams) Service {
	params.validate()

	return &service{
		repo:           params.Repo,
		personService:  params.PersonService,
		dogAPIClient:   params.DogAPIClient,
		holidayService: params.HolidayService,
	}
}

func (p *ServiceParams) validate() {
	err := validator.ValidateStruct(p)
	if err != nil {
		panic(err)
	}
}

func (s *service) ListPacks(ctx context.Context, filters *ListFilters) ([]*Entity, *pagination.Metadata, error) {
	if filters.PageSize == 0 {
		filters.PageSize = 100
	}

	if filters.PageSize > 1000 {
		filters.PageSize = 1000
	}

	packs, metadata, err := s.repo.List(ctx, filters)
	if err != nil {
		return nil, nil, err
	}

	return packs, metadata, nil
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

func (s *service) GetPackByID(ctx context.Context, id string, withEvents bool) (*Entity, error) {
	pack, err := s.repo.GetByID(ctx, id, withEvents)
	if err != nil {
		return nil, err
	}

	if pack == nil {
		return nil, ErrPackNotFound
	}

	return pack, nil
}

func (s *service) UpdatePackStatusByID(ctx context.Context, id string, pack *Entity) (*Entity, error) {
	currentPack, err := s.GetPackByID(ctx, id, false)
	if err != nil {
		return nil, err
	}

	err = currentPack.Status.ValidateChangeStatus(pack.Status)
	if err != nil {
		return nil, err
	}

	currentPack.Status = pack.Status

	if currentPack.Status == StatusDelivered {
		now := time.Now()
		currentPack.DeliveredAt = &now
	}

	err = s.repo.UpdateByID(ctx, id, currentPack)
	if err != nil {
		return nil, err
	}

	return currentPack, nil
}

func (s *service) CancelPackStatusByID(ctx context.Context, id string) (*Entity, error) {
	currentPack, err := s.GetPackByID(ctx, id, false)
	if err != nil {
		return nil, err
	}

	if currentPack.Status != StatusCreated {
		return nil, ErrCannotCancel
	}

	now := time.Now()
	currentPack.Status = StatusCanceled
	currentPack.CanceledAt = &now

	err = s.repo.UpdateByID(ctx, id, currentPack)
	if err != nil {
		return nil, err
	}

	return currentPack, nil
}

func (s *service) setFunFact(ctx context.Context, pack *Entity) {
	funFacts, err := s.dogAPIClient.GetRandomFacts(ctx, 1)
	if err != nil {
		// TODO: implement error handler
		log.Printf("Error getting fun facts: %s. pack: %s", err, pack.ID)
		return
	}

	if len(funFacts) == 0 {
		log.Printf("No fun facts found to pack: %s", pack.ID)
		return
	}

	pack.FunFact = &funFacts[0].Attributes.Body

	err = s.repo.UpdateByID(ctx, pack.ID, pack)
	if err != nil {
		log.Printf("Error updating pack: (%s). pack: %s", err, pack.ID)
	}
}

func (s *service) setIsHoliday(ctx context.Context, pack *Entity) {
	isHoliday, err := s.holidayService.IsHoliday(ctx, pack.EstimatedDeliveryDate)
	if err != nil {
		// TODO: implement error handler
		log.Printf("Error getting holidays: %s. pack: %s", err, pack.ID)
	}

	pack.IsHoliday = &isHoliday

	err = s.repo.UpdateByID(ctx, pack.ID, pack)
	if err != nil {
		// TODO: implement error handler
		log.Printf("Error updating pack: (%s). pack: %s", err, pack.ID)
	}
}
