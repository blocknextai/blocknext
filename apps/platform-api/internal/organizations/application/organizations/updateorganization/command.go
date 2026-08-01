package updateorganization

import (
	"github.com/google/uuid"
)

type UpdateOrganizationCommand struct {
	UserID         uuid.UUID
	OrganizationID uuid.UUID
	Title          string
	Description    *string
}
