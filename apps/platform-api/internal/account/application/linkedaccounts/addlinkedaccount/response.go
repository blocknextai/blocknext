package addlinkedaccount

import (
	"github.com/google/uuid"
)

type AddLinkedAccountResponse struct {
	ID           uuid.UUID `json:"id"`
	AuthProvider string    `json:"authProvider"`
	Identifier   string    `json:"identifier"`
	DisplayName  *string   `json:"displayName"`
	IsPrimary    bool      `json:"isPrimary"`
}
