package workflows

import (
	"strings"
	"time"

	"github.com/blocknextai/go-packages/dag"
	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/google/uuid"
)

type Workflow struct {
	database.BaseEntity

	OrganizationID uuid.UUID
	OwnerID        uuid.UUID
	Title          string
	Description    *string
	IsPinned       bool
	IsActive       bool
	Nodes          []dag.Node
	Edges          []dag.Edge
}

func New(
	organizationID uuid.UUID,
	ownerID uuid.UUID,
	title string,
	description *string,
	isPinned bool,
	isActive bool,
	nodes []dag.Node,
	edges []dag.Edge,
) (*Workflow, error) {
	utcNow := time.Now().UTC()

	workflow := &Workflow{
		ID:             bnuuid.NewV7(),
		CreatedAt:      utcNow,
		UpdatedAt:      utcNow,
		DeletedAt:      nil,
		OrganizationID: organizationID,
		OwnerID:        ownerID,
		Title:          title,
		Description:    description,
		IsPinned:       isPinned,
		IsActive:       isActive,
		Nodes:          nodes,
		Edges:          edges,
	}

	return workflow.validateThenReturn()
}

func (w *Workflow) Update(
	title string,
	description *string,
	isPinned bool,
	isActive bool,
	nodes []dag.Node,
	edges []dag.Edge,
) (*Workflow, error) {
	w.UpdatedAt = time.Now().UTC()

	w.Title = title
	w.Description = description
	w.IsPinned = isPinned
	w.IsActive = isActive
	w.Nodes = nodes
	w.Edges = edges

	return w.validateThenReturn()
}

func (w *Workflow) Delete() (*Workflow, error) {
	utcNow := time.Now().UTC()

	w.UpdatedAt = utcNow
	w.DeletedAt = new(utcNow)
	return w.validateThenReturn()
}

func (w *Workflow) validateThenReturn() (*Workflow, error) {
	if w.OrganizationID == uuid.Nil {
		return nil, ErrOrganizationIDIsRequired
	}

	if w.OwnerID == uuid.Nil {
		return nil, ErrOwnerIDIsRequired
	}

	if strings.TrimSpace(w.Title) == "" {
		return nil, ErrTitleIsRequired
	}

	return w, nil
}
