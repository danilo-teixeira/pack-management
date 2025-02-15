package packevent

import "time"

type (
	Entity struct {
		ID          string
		PackID      string
		Description string
		Location    string
		Date        time.Time
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
)

func (e *Entity) ToModel() *Model {
	if e == nil {
		return nil
	}

	model := &Model{
		ID:          e.ID,
		PackID:      e.PackID,
		Description: e.Description,
		Location:    e.Location,
		Date:        e.Date,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}

	return model
}
