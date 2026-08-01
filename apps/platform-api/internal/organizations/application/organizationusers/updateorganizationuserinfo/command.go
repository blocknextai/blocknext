package updateorganizationuserinfo

import (
	"github.com/google/uuid"
)

type UpdateOrganizationUserInfoCommand struct {
	UserID             uuid.UUID
	OrganizationID     uuid.UUID
	OrganizationUserID uuid.UUID
	Alias              *string
}
