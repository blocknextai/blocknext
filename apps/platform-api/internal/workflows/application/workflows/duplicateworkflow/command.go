package duplicateworkflow

import (
	"github.com/google/uuid"
)

type DuplicateWorkflowCommand struct {
	WorkflowID     uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Title          string
	Description    *string
}
