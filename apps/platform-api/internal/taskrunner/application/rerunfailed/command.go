package rerunfailed

import (
	"github.com/google/uuid"
)

type RerunFailedCommand struct {
	TriggeredByUserID uuid.UUID
	OrganizationID    uuid.UUID
	ID                uuid.UUID
}
