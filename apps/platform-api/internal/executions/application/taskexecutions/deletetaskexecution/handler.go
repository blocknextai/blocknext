package deletetaskexecution

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	executionsApplicationTaskExecutions "github.com/blocknextai/platform-api/internal/executions/application/taskexecutions"
	"github.com/blocknextai/platform-api/internal/executions/domain/taskexecutions"
)

type Handler struct {
	taskExecutionRepository taskexecutions.TaskExecutionRepository
	transactionManager      database.TransactionManager
}

func New(
	taskExecutionRepository taskexecutions.TaskExecutionRepository,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		taskExecutionRepository: taskExecutionRepository,
		transactionManager:      transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, request *DeleteTaskExecutionCommand) (*DeleteTaskExecutionResponse, error) {
	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		taskExecution, err := h.taskExecutionRepository.GetByIDAndOrganizationID(txCtx, request.ID, request.OrganizationID)
		if err != nil {
			return err
		}

		deletedExecution, err := taskExecution.Delete()
		if err != nil {
			return err
		}

		err = h.taskExecutionRepository.Delete(txCtx, deletedExecution)
		if err != nil {
			return executionsApplicationTaskExecutions.ErrFailedToDeleteTaskExecution.WithCause(err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &DeleteTaskExecutionResponse{}, nil
}
