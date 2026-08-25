package workflowresolver

import (
	"context"
	"log/slog"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	workflowsApplicationWorkflows "github.com/blocknextai/platform-api/internal/workflows/application/workflows"
	"github.com/google/uuid"
)

const deletedTitle = "[deleted]"

type Workflow struct {
	ID    uuid.UUID
	Title string
}

type Resolver interface {
	Resolve(
		ctx context.Context,
		executionContext commonDomain.ExecutionContext,
		contextItemID uuid.UUID,
		organizationID uuid.UUID,
	) Workflow
}

type resolver struct {
	workflowService workflowsApplicationWorkflows.WorkflowService
}

func New(workflowService workflowsApplicationWorkflows.WorkflowService) Resolver {
	return &resolver{
		workflowService: workflowService,
	}
}

func (r *resolver) Resolve(
	ctx context.Context,
	executionContext commonDomain.ExecutionContext,
	contextItemID uuid.UUID,
	organizationID uuid.UUID,
) Workflow {
	if executionContext == commonDomain.ExecutionContextWorkflow {
		workflow, err := r.workflowService.GetWorkflow(ctx, organizationID, contextItemID)
		if err == nil {
			return Workflow{ID: workflow.ID, Title: workflow.Title}
		}
		warn(ctx, "failed to resolve workflow for task execution", organizationID, contextItemID, err)
	}

	return Workflow{Title: deletedTitle}
}

func warn(ctx context.Context, message string, organizationID uuid.UUID, contextItemID uuid.UUID, err error) {
	slog.WarnContext(ctx, message,
		"component", "workflow_resolver",
		"organization_id", organizationID,
		"context_item_id", contextItemID,
		"error", err)
}
