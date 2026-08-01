package updatesession

import (
	"github.com/google/uuid"
)

type UpdateSessionCommand struct {
	OrganizationID uuid.UUID
	SessionID      uuid.UUID
	Title          string
}
