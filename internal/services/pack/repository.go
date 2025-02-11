package pack

import (
	"context"
	"pack-management/internal/services/person"
	"time"

	"github.com/uptrace/bun"
)

type (
	Repository interface {
		Create(ctx context.Context, pack *Entity) error
		UpdateByID(ctx context.Context, ID string, pack *Entity) error
		GetByID(ctx context.Context, ID string) (*Entity, error)
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
	}
}
