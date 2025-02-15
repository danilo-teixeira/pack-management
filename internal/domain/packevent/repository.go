package packevent

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type (
	Repository interface {
		Create(ctx context.Context, event *Entity) error
	}

	Model struct {
		bun.BaseModel `bun:"table:pack_event,alias:pack_event"`
		ID            string    `bun:"id,pk"`
		PackID        string    `bun:"pack_id"`
		Description   string    `bun:"description"`
		Location      string    `bun:"location"`
		Date          time.Time `bun:"date"`
		CreatedAt     time.Time `bun:"created_at"`
		UpdatedAt     time.Time `bun:"updated_at"`
	}
)

const (
	idPrefix = "event_"
)

func (m *Model) ToEntity() *Entity {
	if m == nil {
		return nil
	}

	return &Entity{
		ID:          m.ID,
		PackID:      m.PackID,
		Description: m.Description,
		Location:    m.Location,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		Date:        m.Date,
	}
}
