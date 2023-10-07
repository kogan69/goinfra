package utils

import (
	"github.com/gofrs/uuid/v5"
)

func IdOrNil() uuid.UUID {
	id, _ := uuid.NewV4()
	return id
}
