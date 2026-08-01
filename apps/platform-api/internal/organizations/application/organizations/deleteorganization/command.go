package deleteorganization

import (
	"github.com/google/uuid"
)

type DeleteOrganizationCommand struct {
	UserID         uuid.UUID
	OrganizationID uuid.UUID
}
