package createcredential

import (
	"github.com/google/uuid"
)

type CreateCredentialResponse struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Key   string    `json:"key"`
	UIKey string    `json:"uiKey"`
}
