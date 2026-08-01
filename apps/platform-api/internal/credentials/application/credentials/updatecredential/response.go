package updatecredential

import (
	"github.com/google/uuid"
)

type UpdateCredentialResponse struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Key   string    `json:"key"`
	UIKey string    `json:"uiKey"`
}
