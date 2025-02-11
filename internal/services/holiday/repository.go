package holiday

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type (
	Repository interface {
		Create(ctx context.Context, holiday *Entity) error
		BulkCreate(ctx context.Context, holidays []*Entity) error
		GetByDate(ctx context.Context, date string) (*Entity, error)
	}

	Model struct {
		bun.BaseModel `bun:"table:holiday,alias:holiday"`
		ID            string    `bun:"id,pk"`
		Name          string    `bun:"name"`
		Date          string    `bun:"date"`
		CreatedAt     time.Time `bun:"created_at"`
		UpdatedAt     time.Time `bun:"updated_at"`
	}
)

const (
	idPrefix = "holiday_"
)

func (m *Model) ToEntity() *Entity {
	if m == nil {
		return nil
	}

	return &Entity{
		ID:        m.ID,
		Name:      m.Name,
		Date:      m.Date,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
