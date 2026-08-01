package createapikey

import (
	"github.com/google/uuid"
)

type CreateAPIKeyResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Key  string    `json:"key"`
}
