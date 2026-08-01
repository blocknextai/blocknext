package createworkflow

import (
	"github.com/blocknextai/go-packages/dag"
	"github.com/google/uuid"
)

type CreateWorkflowCommand struct {
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Title          string
	Description    string
	Nodes          []dag.Node
	Edges          []dag.Edge
}
