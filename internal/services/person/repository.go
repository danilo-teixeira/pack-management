package person

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type (
	Repository interface {
		Create(ctx context.Context, person *Entity) error
		GetByName(ctx context.Context, name string) (*Entity, error)
	}

	Model struct {
		bun.BaseModel `bun:"table:person,alias:person"`
		ID            string    `bun:"id,pk"`
		Name          string    `bun:"name"`
		CreatedAt     time.Time `bun:"created_at"`
		UpdatedAt     time.Time `bun:"updated_at"`
	}
)

var idPrefix = "person_"

func (m *Model) ToEntity() *Entity {
	if m == nil {
		return nil
	}

	return &Entity{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
