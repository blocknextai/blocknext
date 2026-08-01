package updateorganizationuserinfo

import (
	"github.com/google/uuid"
)

type UpdateOrganizationUserInfoResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organizationId"`
	UserID         uuid.UUID `json:"userId"`
	Role           string    `json:"role"`
	Alias          *string   `json:"alias"`
}
