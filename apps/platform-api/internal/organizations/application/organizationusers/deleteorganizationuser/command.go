package deleteorganizationuser

import (
	"github.com/google/uuid"
)

type DeleteOrganizationUserCommand struct {
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	ForceDelete    bool
}
