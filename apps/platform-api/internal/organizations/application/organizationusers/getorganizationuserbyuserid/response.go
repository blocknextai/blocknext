package getorganizationuserbyuserid

import (
	"github.com/google/uuid"
)

type LinkedAccount struct {
	AuthProvider string  `json:"authProvider"`
	DisplayName  *string `json:"displayName"`
	IsPrimary    bool    `json:"isPrimary"`
}

type GetOrganizationUserByUserIDResponse struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organizationId"`
	UserID         uuid.UUID       `json:"userId"`
	Role           string          `json:"role"`
	Permissions    []string        `json:"permissions"`
	Alias          string          `json:"alias"`
	IsVerified     bool            `json:"isVerified"`
	LinkedAccounts []LinkedAccount `json:"linkedAccounts"`
}
