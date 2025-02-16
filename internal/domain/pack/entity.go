package pack

import (
	"pack-management/internal/domain/person"
	"pack-management/internal/pkg/cerrors"
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
		Events                []*EventEntity
	}

	EventEntity struct {
		ID          string
		PackID      string
		Description string
		Location    string
		Date        time.Time
	}

	Status string
)

var (
	StatusCreated   Status = "CREATED"
	StatusInTransit Status = "IN_TRANSIT"
	StatusDelivered Status = "DELIVERED"
	StatusCanceled  Status = "CANCELED"

	ErrPackNotFound  = cerrors.New("pack not found", "pack_not_found")
	ErrStatusInvalid = cerrors.New("the informed status is invalid", "status_invalid")
	ErrCannotCancel  = cerrors.New("cannot cancel pack is already sent", "cannot_cancel")
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

func (e *EventEntity) ToModel() *EventModel {
	if e == nil {
		return nil
	}

	model := &EventModel{
		ID:          e.ID,
		PackID:      e.PackID,
		Description: e.Description,
		Location:    e.Location,
		Date:        e.Date,
	}

	return model
}

func (s *Status) String() string {
	if s == nil {
		return ""
	}

	return string(*s)
}

func (s *Status) ValidateChangeStatus(newStatus Status) error {
	if s == nil {
		return nil
	}

	if *s == StatusDelivered {
		return ErrStatusInvalid
	}

	if *s == StatusCanceled {
		return ErrStatusInvalid
	}

	if *s == StatusCreated && newStatus != StatusInTransit {
		return ErrStatusInvalid
	}

	if *s == StatusInTransit && newStatus != StatusDelivered && newStatus != StatusCanceled {
		return ErrStatusInvalid
	}

	return nil
}
