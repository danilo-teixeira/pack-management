package uuid

import (
	"github.com/gofrs/uuid/v5"
)

func New() uuid.UUID {
	id, err := uuid.NewV6()
	if err != nil {
		panic(err)
	}

	return id
}
