package taskexecutions

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/apperror"
)

type BulkDeleteTaskExecutionsCommand struct {
	IDs            []uuid.UUID
	OrganizationID uuid.UUID
}

var (
	ErrIDsIsRequired = apperror.Validation("ids is required")
	ErrIDsIsEmpty    = apperror.Validation("ids is empty")
	ErrTooManyIDs    = apperror.Validation("too many ids")
)

const MaxBulkDeleteLimit = 50

func (c *BulkDeleteTaskExecutionsCommand) Validate() error {
	if c.IDs == nil {
		return ErrIDsIsRequired
	}

	if len(c.IDs) == 0 {
		return ErrIDsIsEmpty
	}

	if len(c.IDs) > MaxBulkDeleteLimit {
		return ErrTooManyIDs
	}

	return nil
}

type BulkDeleteTaskExecutionsResponse struct {
	DeletedCount int         `json:"deletedCount"`
	DeletedIDs   []uuid.UUID `json:"deletedIds"`
	FailedIDs    []uuid.UUID `json:"failedIds"`
}

func (s *Service) BulkDeleteTaskExecutions(ctx context.Context, request *BulkDeleteTaskExecutionsCommand) (*BulkDeleteTaskExecutionsResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	var response *BulkDeleteTaskExecutionsResponse

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		validTaskIDs, failedIDs, err := s.validateTaskExecutions(txCtx, request.IDs, request.OrganizationID)
		if err != nil {
			return err
		}

		deletedIDs, deleteFailedIDs := s.performBulkDelete(txCtx, validTaskIDs, request.OrganizationID)

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

func (s *Service) validateTaskExecutions(ctx context.Context, ids []uuid.UUID, organizationID uuid.UUID) ([]uuid.UUID, []uuid.UUID, error) {
	existingTaskExecutions, err := s.taskExecutionRepository.GetAllByIDsAndOrganizationID(ctx, ids, organizationID)
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

func (s *Service) performBulkDelete(ctx context.Context, validTaskIDs []uuid.UUID, organizationID uuid.UUID) ([]uuid.UUID, []uuid.UUID) {
	if len(validTaskIDs) == 0 {
		return nil, nil
	}

	utcNow := time.Now().UTC()
	err := s.taskExecutionRepository.BulkDelete(ctx, validTaskIDs, organizationID, utcNow)
	if err != nil {
		return nil, validTaskIDs
	}

	return validTaskIDs, nil
}
