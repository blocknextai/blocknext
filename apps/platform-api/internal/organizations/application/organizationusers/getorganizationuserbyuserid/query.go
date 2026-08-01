package getorganizationuserbyuserid

import (
	"github.com/google/uuid"
)

type GetOrganizationUserByUserIDQuery struct {
	OrganizationID uuid.UUID
	UserID         uuid.UUID
}
