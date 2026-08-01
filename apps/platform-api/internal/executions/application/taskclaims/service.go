package taskclaims

import (
	"context"
	"errors"
	"time"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/executions/domain/taskclaims"
	"github.com/google/uuid"
)

type TaskClaimService interface {
	Create(ctx context.Context, taskExecutionID uuid.UUID) error
	ClaimByID(ctx context.Context, id uuid.UUID, workerID string) (bool, error)
	ReleaseClaim(ctx context.Context, id uuid.UUID, workerID string) error
	ReleaseStaleClaims(ctx context.Context, staleAfter time.Duration) ([]uuid.UUID, error)
	Heartbeat(ctx context.Context, id uuid.UUID, workerID string) error
	IncrementRetryCount(ctx context.Context, id uuid.UUID) (int, error)
	ResetRetryCount(ctx context.Context, id uuid.UUID) error
}

type taskClaimService struct {
	taskClaimRepository taskclaims.TaskClaimRepository
	transactionManager  database.TransactionManager
}

func NewTaskClaimService(
	taskClaimRepository taskclaims.TaskClaimRepository,
	transactionManager database.TransactionManager,
) TaskClaimService {
	return &taskClaimService{
		taskClaimRepository: taskClaimRepository,
		transactionManager:  transactionManager,
	}
}

func (s *taskClaimService) Create(ctx context.Context, taskExecutionID uuid.UUID) error {
	taskClaim, err := taskclaims.New(taskExecutionID)
	if err != nil {
		return err
	}
	return s.taskClaimRepository.Create(ctx, taskClaim)
}

func (s *taskClaimService) ClaimByID(ctx context.Context, id uuid.UUID, workerID string) (bool, error) {
	var claimed bool
	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		taskClaim, err := s.taskClaimRepository.GetByTaskExecutionIDForUpdate(txCtx, id)
		if err != nil {
			return err
		}

		updated, err := taskClaim.Claim(workerID)
		if errors.Is(err, taskclaims.ErrTaskClaimAlreadyClaimed) {
			claimed = false
			return nil
		}
		if err != nil {
			return err
		}

		if err := s.taskClaimRepository.Update(txCtx, updated); err != nil {
			return err
		}
		claimed = true
		return nil
	})
	return claimed, err
}

func (s *taskClaimService) ReleaseClaim(ctx context.Context, id uuid.UUID, workerID string) error {
	return s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		taskClaim, err := s.taskClaimRepository.GetByTaskExecutionIDForUpdate(txCtx, id)
		if err != nil {
			return err
		}

		updated, err := taskClaim.Release(workerID)
		if errors.Is(err, taskclaims.ErrTaskClaimNotOwner) {
			return nil
		}
		if err != nil {
			return err
		}

		return s.taskClaimRepository.Update(txCtx, updated)
	})
}

func (s *taskClaimService) ReleaseStaleClaims(ctx context.Context, staleAfter time.Duration) ([]uuid.UUID, error) {
	released := make([]uuid.UUID, 0)
	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		staleClaims, err := s.taskClaimRepository.GetAllStale(txCtx, staleAfter)
		if err != nil {
			return err
		}

		for _, staleClaim := range staleClaims {
			updated, err := staleClaim.ForceRelease()
			if err != nil {
				return err
			}
			if err := s.taskClaimRepository.Update(txCtx, updated); err != nil {
				return err
			}
			released = append(released, staleClaim.TaskExecutionID)
		}
		return nil
	})
	return released, err
}

func (s *taskClaimService) Heartbeat(ctx context.Context, id uuid.UUID, workerID string) error {
	return s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		taskClaim, err := s.taskClaimRepository.GetByTaskExecutionIDForUpdate(txCtx, id)
		if err != nil {
			return err
		}

		updated, err := taskClaim.Heartbeat(workerID)
		if errors.Is(err, taskclaims.ErrTaskClaimNotOwner) {
			return nil
		}
		if err != nil {
			return err
		}

		return s.taskClaimRepository.Update(txCtx, updated)
	})
}

func (s *taskClaimService) IncrementRetryCount(ctx context.Context, id uuid.UUID) (int, error) {
	var newCount int
	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		taskClaim, err := s.taskClaimRepository.GetByTaskExecutionIDForUpdate(txCtx, id)
		if err != nil {
			return err
		}

		updated, err := taskClaim.IncrementRetryCount()
		if err != nil {
			return err
		}

		if err := s.taskClaimRepository.Update(txCtx, updated); err != nil {
			return err
		}
		newCount = updated.RetryCount
		return nil
	})
	return newCount, err
}

func (s *taskClaimService) ResetRetryCount(ctx context.Context, id uuid.UUID) error {
	return s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		taskClaim, err := s.taskClaimRepository.GetByTaskExecutionIDForUpdate(txCtx, id)
		if err != nil {
			return err
		}

		updated, err := taskClaim.ResetRetryCount()
		if err != nil {
			return err
		}

		return s.taskClaimRepository.Update(txCtx, updated)
	})
}
