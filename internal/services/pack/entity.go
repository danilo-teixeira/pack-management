package pack

import (
	"pack-management/internal/services/person"
	"time"
)

type (
	Entity struct {
		ID                    string
		Description           string
		FunFact               *string
		IsHoliday             *bool
		Status                Status
		Receiver              *person.Entity
		Sender                *person.Entity
		EstimatedDeliveryDate string
		DeliveredAt           *time.Time
		CanceledAt            *time.Time
		CreatedAt             time.Time
		UpdatedAt             time.Time
	}

	Status string
)

var (
	StatusCreated   Status = "CREATED"
	StatusInTransit Status = "IN_TRANSIT"
	StatusDelivered Status = "DELIVERED"
	StatusCanceled  Status = "CANCELED"
)

func (e *Entity) ToModel() *Model {
	if e == nil {
		return nil
	}

	estimatedDeliveryDate := time.Time{}
	if e.EstimatedDeliveryDate != "" {
		withTime, _ := time.Parse(time.DateOnly, e.EstimatedDeliveryDate)
		estimatedDeliveryDate = time.Date(withTime.Year(), withTime.Month(), withTime.Day(), 0, 0, 0, 0, time.Local)
	}

	model := &Model{
		ID:                    e.ID,
		Description:           e.Description,
		FunFact:               e.FunFact,
		IsHoliday:             e.IsHoliday,
		Status:                e.Status,
		EstimatedDeliveryDate: estimatedDeliveryDate,
		DeliveredAt:           e.DeliveredAt,
		CanceledAt:            e.CanceledAt,
		CreatedAt:             e.CreatedAt,
		UpdatedAt:             e.UpdatedAt,
	}

	if e.Receiver != nil {
		model.ReceiverID = &e.Receiver.ID
	}

	if e.Sender != nil {
		model.SenderID = &e.Sender.ID
	}

	return model
}

func (s *Status) String() string {
	if s == nil {
		return ""
	}

	return string(*s)
}
