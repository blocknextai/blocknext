package gettoolinvocationbyid

import (
	"github.com/google/uuid"
)

type GetToolInvocationByIDQuery struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
}
