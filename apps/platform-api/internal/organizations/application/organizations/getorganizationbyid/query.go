package getorganizationbyid

import (
	"github.com/google/uuid"
)

type GetOrganizationByIDQuery struct {
	UserID         uuid.UUID
	OrganizationID uuid.UUID
}
