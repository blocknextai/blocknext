package getworkflowforrun

import (
	"github.com/google/uuid"
)

type GetWorkflowForRunQuery struct {
	OrganizationID uuid.UUID
	WorkflowID     uuid.UUID
}
