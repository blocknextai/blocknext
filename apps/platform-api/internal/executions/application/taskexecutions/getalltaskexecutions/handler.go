package getalltaskexecutions

import (
	"context"
	"log/slog"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/blocknextai/platform-api/internal/executions/domain/taskexecutions"
	workflowsApplicationWorkflows "github.com/blocknextai/platform-api/internal/workflows/application/workflows"
	"github.com/google/uuid"
)

type Handler struct {
	taskExecutionRepository taskexecutions.TaskExecutionRepository
	workflowService         workflowsApplicationWorkflows.WorkflowService
}

func New(
	taskExecutionRepository taskexecutions.TaskExecutionRepository,
	workflowService workflowsApplicationWorkflows.WorkflowService,
) *Handler {
	return &Handler{
		taskExecutionRepository: taskExecutionRepository,
		workflowService:         workflowService,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetAllTaskExecutionsQuery) (*GetAllTaskExecutionsResponse, error) {
	taskExecutions, totalCount, err := h.taskExecutionRepository.GetAllByOrganizationID(
		ctx,
		request.OrganizationID,
		request.Search.Query,
		request.Pagination.Offset,
		request.Pagination.Limit,
	)
	if err != nil {
		return nil, err
	}

	workflowsByID := make(map[uuid.UUID]Workflow, len(taskExecutions))
	for _, taskExecution := range taskExecutions {
		workflowsByID[taskExecution.ID] = h.resolveWorkflow(ctx, taskExecution.ExecutionContext, taskExecution.ContextItemID, taskExecution.OrganizationID)
	}

	return &GetAllTaskExecutionsResponse{
		Items:      MapGetAllTaskExecutionsQueryToGetAllTaskExecutionsResponse(taskExecutions, workflowsByID),
		TotalCount: totalCount,
	}, nil
}

func (h *Handler) resolveWorkflow(
	ctx context.Context,
	executionContext commonDomain.ExecutionContext,
	contextItemID uuid.UUID,
	organizationID uuid.UUID,
) Workflow {
	switch executionContext {
	case commonDomain.ExecutionContextWorkflow:
		workflow, err := h.workflowService.GetWorkflow(ctx, organizationID, contextItemID)
		if err != nil {
			slog.WarnContext(ctx, "Failed to resolve workflow info for task execution",
				"component", "getalltaskexecutions",
				"organization_id", organizationID,
				"context_item_id", contextItemID,
				"error", err)
			return Workflow{Title: "[deleted]"}
		}
		return Workflow{
			ID:    workflow.ID,
			Title: workflow.Title,
		}
	default:
		return Workflow{Title: "[deleted]"}
	}
}
