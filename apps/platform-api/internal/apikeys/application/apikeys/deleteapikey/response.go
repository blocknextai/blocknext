package deleteapikey

import (
	"github.com/google/uuid"
)

type DeleteAPIKeyResponse struct {
	ID uuid.UUID `json:"id"`
}
