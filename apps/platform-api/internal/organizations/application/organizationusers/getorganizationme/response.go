package getorganizationme

import (
	"github.com/google/uuid"
)

type GetOrganizationMeResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organizationId"`
	UserID         uuid.UUID `json:"userId"`
	Alias          string    `json:"alias"`
	Role           string    `json:"role"`
	Permissions    []string  `json:"permissions"`
}
