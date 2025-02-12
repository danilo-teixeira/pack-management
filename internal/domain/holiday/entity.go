package holiday

import (
	"time"
)

type (
	Entity struct {
		ID        string
		Name      string
		Date      string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

func (e *Entity) ToModel() *Model {
	if e == nil {
		return nil
	}

	model := &Model{
		ID:        e.ID,
		Name:      e.Name,
		Date:      e.Date,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}

	return model
}
