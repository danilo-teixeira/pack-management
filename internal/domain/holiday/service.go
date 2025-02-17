package holiday

import (
	"context"
	naegerdateapi "pack-management/internal/pkg/http/nagerdateapi"
	"pack-management/internal/pkg/validator"
	"slices"
)

type (
	Service interface {
		IsHoliday(ctx context.Context, date string) (bool, error)
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

func (s *service) IsHoliday(ctx context.Context, date string) (bool, error) {
	year := date[:4]

	holidays, err := s.repo.ListByYear(ctx, year)
	if err != nil {
		return false, err
	}

	if len(holidays) <= 0 {
		holidays, err = s.getHolidaysFromProvider(ctx, year)
		if err != nil {
			return false, err
		}
	}

	isHoliday := slices.ContainsFunc(holidays, func(holiday *Entity) bool {
		return holiday != nil && holiday.Date == date
	})

	return isHoliday, nil
}

func (s *service) getHolidaysFromProvider(ctx context.Context, year string) ([]*Entity, error) {
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

	return holidays, nil
}
