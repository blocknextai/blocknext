package regenerateapikey

import (
	"github.com/google/uuid"
)

type RegenerateAPIKeyResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Key  string    `json:"key"`
}
