package updateorganizationuserrole

import (
	"github.com/google/uuid"
)

type UpdateOrganizationUserRoleResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organizationId"`
	UserID         uuid.UUID `json:"userId"`
	Role           string    `json:"role"`
}
