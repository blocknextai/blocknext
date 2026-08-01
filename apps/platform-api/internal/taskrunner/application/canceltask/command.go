package canceltask

import (
	"github.com/google/uuid"
)

type CancelTaskCommand struct {
	TriggeredByUserID uuid.UUID
	OrganizationID    uuid.UUID
	ID                uuid.UUID
}
