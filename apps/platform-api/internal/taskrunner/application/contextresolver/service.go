package contextresolver

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/dag"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/blocknextai/platform-api/internal/workflows/application/workflows"
	"github.com/google/uuid"
)

type ResolvedContext struct {
	ContextType string
	ItemID      uuid.UUID
	ItemType    string
	Nodes       []dag.Node
	Edges       []dag.Edge
}

type ContextResolver interface {
	ResolveContext(ctx context.Context, contextType commonDomain.ExecutionContext, itemID uuid.UUID, organizationID uuid.UUID) (*ResolvedContext, error)
	GetWorkflowData(ctx context.Context, contextType commonDomain.ExecutionContext, itemID uuid.UUID, organizationID uuid.UUID) ([]dag.Node, []dag.Edge, error)
}

type contextResolver struct {
	workflowService workflows.WorkflowService
}

func NewContextResolver(
	workflowService workflows.WorkflowService,
) ContextResolver {
	return &contextResolver{
		workflowService: workflowService,
	}
}

func (r *contextResolver) ResolveContext(ctx context.Context, contextType commonDomain.ExecutionContext, itemID uuid.UUID, organizationID uuid.UUID) (*ResolvedContext, error) {
	contextData, err := r.getContextData(ctx, contextType, itemID, organizationID)
	if err != nil {
		return nil, err
	}

	return &ResolvedContext{
		ContextType: contextType.String(),
		ItemID:      itemID,
		ItemType:    contextData.ItemType,
		Nodes:       contextData.Nodes,
		Edges:       contextData.Edges,
	}, nil
}

func (r *contextResolver) GetWorkflowData(ctx context.Context, contextType commonDomain.ExecutionContext, itemID uuid.UUID, organizationID uuid.UUID) ([]dag.Node, []dag.Edge, error) {
	contextData, err := r.getContextData(ctx, contextType, itemID, organizationID)
	if err != nil {
		return nil, nil, err
	}

	return contextData.Nodes, contextData.Edges, nil
}

type contextData struct {
	ItemType   string
	WorkflowID uuid.UUID
	Nodes      []dag.Node
	Edges      []dag.Edge
}

func (r *contextResolver) getContextData(ctx context.Context, contextType commonDomain.ExecutionContext, itemID uuid.UUID, organizationID uuid.UUID) (*contextData, error) {
	switch contextType {
	case commonDomain.ExecutionContextWorkflow:
		workflow, err := r.workflowService.GetWorkflow(ctx, organizationID, itemID)
		if err != nil {
			return nil, err
		}

		return &contextData{
			ItemType:   commonDomain.ExecutionContextWorkflow.String(),
			WorkflowID: workflow.ID,
			Nodes:      workflow.Nodes,
			Edges:      workflow.Edges,
		}, nil

	default:
		return nil, apperror.Validation("unsupported context type: " + contextType.String())
	}
}
