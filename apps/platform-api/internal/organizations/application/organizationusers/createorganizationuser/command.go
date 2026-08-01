package createorganizationuser

import (
	"github.com/google/uuid"
)

type CreateOrganizationUserCommand struct {
	OrganizationID uuid.UUID
	Identifier     string
	Alias          string
	Role           string
}
