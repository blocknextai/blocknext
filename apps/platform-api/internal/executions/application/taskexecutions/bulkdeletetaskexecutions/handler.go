package bulkdeletetaskexecutions

import (
	"context"
	"time"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/executions/domain/taskexecutions"
	"github.com/google/uuid"
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

func (h *Handler) Handle(ctx context.Context, request *BulkDeleteTaskExecutionsCommand) (*BulkDeleteTaskExecutionsResponse, error) {
	var response *BulkDeleteTaskExecutionsResponse

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		validTaskIDs, failedIDs, err := h.validateTaskExecutions(txCtx, request.IDs, request.OrganizationID)
		if err != nil {
			return err
		}

		deletedIDs, deleteFailedIDs := h.performBulkDelete(txCtx, validTaskIDs, request.OrganizationID)

		allFailedIDs := make([]uuid.UUID, 0, len(failedIDs)+len(deleteFailedIDs))
		allFailedIDs = append(allFailedIDs, failedIDs...)
		allFailedIDs = append(allFailedIDs, deleteFailedIDs...)

		response = &BulkDeleteTaskExecutionsResponse{
			DeletedCount: len(deletedIDs),
			DeletedIDs:   deletedIDs,
			FailedIDs:    allFailedIDs,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (h *Handler) validateTaskExecutions(ctx context.Context, ids []uuid.UUID, organizationID uuid.UUID) ([]uuid.UUID, []uuid.UUID, error) {
	existingTaskExecutions, err := h.taskExecutionRepository.GetAllByIDsAndOrganizationID(ctx, ids, organizationID)
	if err != nil {
		return nil, nil, err
	}

	existingTaskExecutionIDs := make(map[uuid.UUID]bool, len(existingTaskExecutions))
	for _, taskExecution := range existingTaskExecutions {
		existingTaskExecutionIDs[taskExecution.ID] = true
	}

	validIDs := make([]uuid.UUID, 0, len(existingTaskExecutions))
	failedIDs := make([]uuid.UUID, 0, len(ids)-len(existingTaskExecutions))

	for _, id := range ids {
		if !existingTaskExecutionIDs[id] {
			failedIDs = append(failedIDs, id)
			continue
		}
		validIDs = append(validIDs, id)
	}

	return validIDs, failedIDs, nil
}

func (h *Handler) performBulkDelete(ctx context.Context, validTaskIDs []uuid.UUID, organizationID uuid.UUID) ([]uuid.UUID, []uuid.UUID) {
	if len(validTaskIDs) == 0 {
		return nil, nil
	}

	utcNow := time.Now().UTC()
	err := h.taskExecutionRepository.BulkDelete(ctx, validTaskIDs, organizationID, utcNow)
	if err != nil {
		return nil, validTaskIDs
	}

	return validTaskIDs, nil
}
