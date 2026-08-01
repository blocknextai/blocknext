package deleteworkflow

import (
	"github.com/google/uuid"
)

type DeleteWorkflowCommand struct {
	OrganizationID uuid.UUID
	WorkflowID     uuid.UUID
}
