package updateorganizationuserrole

import (
	"github.com/google/uuid"
)

type UpdateOrganizationUserRoleCommand struct {
	UserID             uuid.UUID
	OrganizationID     uuid.UUID
	OrganizationUserID uuid.UUID
	Role               string
}
