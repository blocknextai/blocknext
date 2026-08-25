package getalltaskexecutions

import (
	"context"

	"github.com/blocknextai/platform-api/internal/executions/application/taskexecutions/workflowresolver"
	"github.com/blocknextai/platform-api/internal/executions/domain/taskexecutions"
	"github.com/google/uuid"
)

type Handler struct {
	taskExecutionRepository taskexecutions.TaskExecutionRepository
	workflowResolver        workflowresolver.Resolver
}

func New(
	taskExecutionRepository taskexecutions.TaskExecutionRepository,
	workflowResolver workflowresolver.Resolver,
) *Handler {
	return &Handler{
		taskExecutionRepository: taskExecutionRepository,
		workflowResolver:        workflowResolver,
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
		resolved := h.workflowResolver.Resolve(ctx, taskExecution.ExecutionContext, taskExecution.ContextItemID, taskExecution.OrganizationID)
		workflowsByID[taskExecution.ID] = Workflow{ID: resolved.ID, Title: resolved.Title}
	}

	return &GetAllTaskExecutionsResponse{
		Items:      MapGetAllTaskExecutionsQueryToGetAllTaskExecutionsResponse(taskExecutions, workflowsByID),
		TotalCount: totalCount,
	}, nil
}
