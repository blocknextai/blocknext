package getworkflowbyid

import (
	"github.com/google/uuid"
)

type GetWorkflowByIDQuery struct {
	OrganizationID uuid.UUID
	WorkflowID     uuid.UUID
}
