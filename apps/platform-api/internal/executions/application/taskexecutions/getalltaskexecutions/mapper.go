package getalltaskexecutions

import (
	executionsDomainTaskexecutions "github.com/blocknextai/platform-api/internal/executions/domain/taskexecutions"
	"github.com/google/uuid"
)

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
