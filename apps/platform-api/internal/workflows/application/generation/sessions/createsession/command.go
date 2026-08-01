package createsession

import (
	"github.com/google/uuid"
)

type CreateSessionCommand struct {
	OrganizationID *uuid.UUID
	UserID         *uuid.UUID
	Title          string
}
