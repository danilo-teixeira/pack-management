package packevent

type (
	EventJSON struct {
		ID          string `json:"id"`
		PackID      string `json:"pack_id"`
		Description string `json:"description"`
		Date        string `json:"date"`
		CreatedAt   string `json:"created_at"`
	}
)
