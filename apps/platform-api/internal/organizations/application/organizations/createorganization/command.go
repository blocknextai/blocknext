package createorganization

import (
	"github.com/google/uuid"
)

type CreateOrganizationCommand struct {
	UserID      uuid.UUID
	Title       string
	Description *string
}
