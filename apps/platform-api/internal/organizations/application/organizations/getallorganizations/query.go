package getallorganizations

import (
	"github.com/google/uuid"
)

type GetAllOrganizationsQuery struct {
	UserID uuid.UUID
}
