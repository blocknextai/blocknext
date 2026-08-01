package rerunall

import (
	"github.com/google/uuid"
)

type RerunAllCommand struct {
	TriggeredByUserID uuid.UUID
	OrganizationID    uuid.UUID
	ID                uuid.UUID
}
