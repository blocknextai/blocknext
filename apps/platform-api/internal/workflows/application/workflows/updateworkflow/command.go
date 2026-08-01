package updateworkflow

import (
	"github.com/blocknextai/go-packages/dag"
	"github.com/google/uuid"
)

type UpdateWorkflowCommand struct {
	OrganizationID uuid.UUID
	WorkflowID     uuid.UUID
	Title          *string
	Description    *string
	IsPinned       *bool
	IsActive       *bool
	Nodes          []dag.Node
	Edges          []dag.Edge
}
