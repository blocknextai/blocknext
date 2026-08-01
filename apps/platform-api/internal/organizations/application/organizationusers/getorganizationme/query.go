package getorganizationme

import (
	"github.com/google/uuid"
)

type GetOrganizationMeQuery struct {
	OrganizationID uuid.UUID
	UserID         uuid.UUID
}
