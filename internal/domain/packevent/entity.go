package packevent

import "time"

type (
	Entity struct {
		ID          string
		PackID      string
		Description string
		Date        time.Time
		CreatedAt   time.Time
	}
)
