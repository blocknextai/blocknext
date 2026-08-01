package deletetaskexecution

import (
	"github.com/google/uuid"
)

type DeleteTaskExecutionCommand struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
}
