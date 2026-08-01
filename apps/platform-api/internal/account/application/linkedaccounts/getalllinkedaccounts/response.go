package getalllinkedaccounts

import (
	"github.com/google/uuid"
)

type LinkedAccount struct {
	ID           uuid.UUID `json:"id"`
	AuthProvider string    `json:"authProvider"`
	Identifier   string    `json:"identifier"`
	DisplayName  *string   `json:"displayName"`
	IsPrimary    bool      `json:"isPrimary"`
	IsVerified   bool      `json:"isVerified"`
}

type GetAllLinkedAccountsResponse = []LinkedAccount
