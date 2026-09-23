package taskexecutions

import (
	"context"

	"github.com/google/uuid"

	executionsApplicationTaskExecutions "github.com/blocknextai/platform-api/internal/modules/executions/internal/application/taskexecutions"
)

type DeleteTaskExecutionCommand struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
}

type DeleteTaskExecutionResponse struct{}

func (s *Service) DeleteTaskExecution(ctx context.Context, request *DeleteTaskExecutionCommand) (*DeleteTaskExecutionResponse, error) {
	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		taskExecution, err := s.taskExecutionRepository.GetByIDAndOrganizationID(txCtx, request.ID, request.OrganizationID)
		if err != nil {
			return err
		}

		deletedExecution, err := taskExecution.Delete()
		if err != nil {
			return err
		}

		err = s.taskExecutionRepository.Delete(txCtx, deletedExecution)
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
