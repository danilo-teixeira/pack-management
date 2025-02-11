package holiday

import (
	"context"
	naegerdateapi "pack-management/internal/pkg/http/nagerdateapi"
	"pack-management/internal/pkg/validator"
)

type (
	Service interface {
		GetByDate(ctx context.Context, date string) (*Entity, error)
	}

	service struct {
		repo   Repository
		client naegerdateapi.Client
	}

	ServiceParams struct {
		Repo   Repository           `validate:"required"`
		Client naegerdateapi.Client `validate:"required"`
	}
)

const (
	countryCode = "BR"
)

func NewService(params *ServiceParams) Service {
	params.validate()

	return &service{
		repo:   params.Repo,
		client: params.Client,
	}
}

func (p *ServiceParams) validate() {
	err := validator.ValidateStruct(p)
	if err != nil {
		panic(err)
	}
}

func (s *service) GetByDate(ctx context.Context, date string) (*Entity, error) {
	holiday, err := s.repo.GetByDate(ctx, date)
	if err != nil {
		return nil, err
	}

	if holiday != nil {
		return holiday, nil
	}

	year := date[:4]

	holidayResponse, err := s.client.GetHolidays(ctx, countryCode, year)
	if err != nil {
		return nil, err
	}

	holidays := make([]*Entity, 0, len(holidayResponse))
	for _, holiday := range holidayResponse {
		holidayEntity := &Entity{
			Name: holiday.Name,
			Date: holiday.Date,
		}

		holidays = append(holidays, holidayEntity)
	}

	err = s.repo.BulkCreate(ctx, holidays)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByDate(ctx, date)
}
