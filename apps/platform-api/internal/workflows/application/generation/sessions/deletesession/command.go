package deletesession

import (
	"github.com/google/uuid"
)

type DeleteSessionCommand struct {
	OrganizationID uuid.UUID
	SessionID      uuid.UUID
}
