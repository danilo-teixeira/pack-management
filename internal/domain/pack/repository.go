package pack

import (
	"context"
	"pack-management/internal/domain/person"
	"pack-management/internal/pkg/pagination"
	"time"

	"github.com/uptrace/bun"
)

type (
	Repository interface {
		Create(ctx context.Context, pack *Entity) error
		List(ctx context.Context, filters *ListFilters) ([]*Entity, *pagination.Metadata, error)
		UpdateByID(ctx context.Context, ID string, pack *Entity) error
		UpdateFunFactByID(ctx context.Context, ID string, funFact string) error
		UpdateIsHolidayByID(ctx context.Context, ID string, isHoliday bool) error
		GetByID(ctx context.Context, ID string, withEvents bool) (*Entity, error)
	}

	Model struct {
		bun.BaseModel         `bun:"table:pack,alias:pack"`
		ID                    string        `bun:"id,pk"`
		Description           string        `bun:"description"`
		FunFact               *string       `bun:"fun_fact"`
		IsHoliday             *bool         `bun:"is_holiday"`
		Status                Status        `bun:"status"`
		EstimatedDeliveryDate time.Time     `bun:"estimated_delivery_date"`
		DeliveredAt           *time.Time    `bun:"delivered_at"`
		CanceledAt            *time.Time    `bun:"canceled_at"`
		CreatedAt             time.Time     `bun:"created_at"`
		UpdatedAt             time.Time     `bun:"updated_at"`
		ReceiverID            *string       `bun:"receiver_id"`
		Receiver              *person.Model `bun:"rel:belongs-to"`
		SenderID              *string       `bun:"sender_id"`
		Sender                *person.Model `bun:"rel:belongs-to"`
		Events                []*EventModel `bun:"rel:has-many,join:id=pack_id"`
	}

	EventModel struct {
		bun.BaseModel `bun:"table:pack_event,alias:pack_event"`
		ID            string    `bun:"id,pk"`
		PackID        string    `bun:"pack_id"`
		Description   string    `bun:"description"`
		Location      string    `bun:"location"`
		Date          time.Time `bun:"date"`
	}
)

const (
	idPrefix = "pack_"
)

func (m *Model) ToEntity() *Entity {
	if m == nil {
		return nil
	}

	estimatedDeliveryDate := m.EstimatedDeliveryDate.Format(time.DateOnly)

	events := make([]*EventEntity, len(m.Events))
	for i, event := range m.Events {
		events[i] = event.ToEntity()
	}

	return &Entity{
		ID:                    m.ID,
		Description:           m.Description,
		FunFact:               m.FunFact,
		IsHoliday:             m.IsHoliday,
		Status:                m.Status,
		EstimatedDeliveryDate: estimatedDeliveryDate,
		DeliveredAt:           m.DeliveredAt,
		CanceledAt:            m.CanceledAt,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
		Receiver:              m.Receiver.ToEntity(),
		Sender:                m.Sender.ToEntity(),
		Events:                events,
	}
}

func (m *EventModel) ToEntity() *EventEntity {
	if m == nil {
		return nil
	}

	return &EventEntity{
		ID:          m.ID,
		PackID:      m.PackID,
		Description: m.Description,
		Location:    m.Location,
		Date:        m.Date,
	}
}
