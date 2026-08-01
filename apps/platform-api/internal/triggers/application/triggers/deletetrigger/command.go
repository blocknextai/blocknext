package deletetrigger

import (
	"github.com/google/uuid"
)

type DeleteTriggerCommand struct {
	OrganizationID uuid.UUID
	TriggerID      uuid.UUID
}
