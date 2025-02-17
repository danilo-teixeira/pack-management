package packevent

import (
	"context"
	"log"
	"pack-management/internal/domain/pack"
	"pack-management/internal/pkg/validator"
)

type (
	Service interface {
		EnqueueEvent(ctx context.Context, event *Entity)
	}

	service struct {
		repo        Repository
		packService pack.Service
		eventsQueue chan *Entity
	}

	ServiceParams struct {
		Repo        Repository   `validate:"required"`
		PackService pack.Service `validate:"required"`
	}
)

const eventsProcessBuffer = 1000

func NewService(ctx context.Context, params *ServiceParams) Service {
	params.validate()

	src := &service{
		repo:        params.Repo,
		packService: params.PackService,
	}

	src.eventsQueue = make(chan *Entity, eventsProcessBuffer)

	go src.processEventsWorker(ctx)

	return src
}

func (p *ServiceParams) validate() {
	err := validator.ValidateStruct(p)
	if err != nil {
		panic(err)
	}
}

func (s *service) EnqueueEvent(ctx context.Context, event *Entity) {
	s.eventsQueue <- event
}

func (s *service) processEventsWorker(ctx context.Context) {
	for {
		select {
		case event := <-s.eventsQueue:
			err := s.createEvent(ctx, event)
			if err != nil {
				log.Printf("Error creating event: %v", err)
				log.Printf("Requeuing event: %v", event)
				s.EnqueueEvent(ctx, event)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *service) createEvent(ctx context.Context, event *Entity) error {
	_, err := s.packService.GetPackByID(ctx, event.PackID, false)
	if err != nil {
		return err
	}

	err = s.repo.Create(ctx, event)
	if err != nil {
		return err
	}

	return nil
}
