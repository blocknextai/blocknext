package taskexecutions

import (
	"context"
	"time"

	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	executionsDomainTaskexecutions "github.com/blocknextai/platform-api/internal/modules/executions/internal/domain/taskexecutions"
)

type GetAllTaskExecutionsQuery struct {
	OrganizationID uuid.UUID
	Search         resultPkg.SearchRequest
	Pagination     resultPkg.PaginationRequest
}

type TaskExecutionResponse struct {
	ID            uuid.UUID  `json:"id,omitempty"`
	Workflow      Workflow   `json:"workflow"`
	Status        string     `json:"status,omitempty"`
	ExecutionType string     `json:"executionType,omitempty"`
	ErrorMessage  *string    `json:"errorMessage,omitempty"`
	StartedAt     *time.Time `json:"startedAt,omitempty"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
}

type Workflow struct {
	ID    uuid.UUID `json:"id,omitempty"`
	Title string    `json:"title,omitempty"`
}

type GetAllTaskExecutionsResponse struct {
	Items      []*TaskExecutionResponse `json:"items"`
	TotalCount int64                    `json:"totalCount"`
}

func (s *Service) GetAllTaskExecutions(ctx context.Context, request *GetAllTaskExecutionsQuery) (*GetAllTaskExecutionsResponse, error) {
	taskExecutions, totalCount, err := s.taskExecutionRepository.GetAllByOrganizationID(
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
		resolved := s.workflowResolver.Resolve(ctx, taskExecution.ExecutionContext, taskExecution.ContextItemID, taskExecution.OrganizationID)
		workflowsByID[taskExecution.ID] = Workflow{ID: resolved.ID, Title: resolved.Title}
	}

	return &GetAllTaskExecutionsResponse{
		Items:      MapGetAllTaskExecutionsQueryToGetAllTaskExecutionsResponse(taskExecutions, workflowsByID),
		TotalCount: totalCount,
	}, nil
}

func MapGetAllTaskExecutionsQueryToGetAllTaskExecutionsResponse(
	taskExecutions []*executionsDomainTaskexecutions.TaskExecution,
	workflowsByID map[uuid.UUID]Workflow,
) []*TaskExecutionResponse {
	items := make([]*TaskExecutionResponse, 0, len(taskExecutions))
	for _, taskExecution := range taskExecutions {
		items = append(items, &TaskExecutionResponse{
			ID:            taskExecution.ID,
			Workflow:      workflowsByID[taskExecution.ID],
			ExecutionType: taskExecution.ExecutionType.String(),
			Status:        taskExecution.Status,
			ErrorMessage:  taskExecution.ErrorMessage,
			StartedAt:     taskExecution.StartedAt,
			CompletedAt:   taskExecution.CompletedAt,
		})
	}

	return items
}
