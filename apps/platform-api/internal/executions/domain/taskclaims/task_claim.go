package taskclaims

import (
	"strings"
	"time"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/google/uuid"
)

type TaskClaim struct {
	database.BaseEntity

	TaskExecutionID uuid.UUID
	ClaimedBy       *string
	ClaimedAt       *time.Time
	RetryCount      int
}

func New(taskExecutionID uuid.UUID) (*TaskClaim, error) {
	utcNow := time.Now().UTC()

	taskClaim := &TaskClaim{
		ID:              bnuuid.NewV7(),
		CreatedAt:       utcNow,
		UpdatedAt:       utcNow,
		DeletedAt:       nil,
		TaskExecutionID: taskExecutionID,
		RetryCount:      0,
	}

	return taskClaim.validateThenReturn()
}

func (c *TaskClaim) Claim(workerID string) (*TaskClaim, error) {
	if strings.TrimSpace(workerID) == "" {
		return nil, ErrWorkerIDIsRequired
	}
	if c.IsClaimed() && *c.ClaimedBy != workerID {
		return nil, ErrTaskClaimAlreadyClaimed
	}

	now := time.Now().UTC()

	c.UpdatedAt = now
	c.ClaimedBy = &workerID
	c.ClaimedAt = &now

	return c.validateThenReturn()
}

func (c *TaskClaim) Release(workerID string) (*TaskClaim, error) {
	if strings.TrimSpace(workerID) == "" {
		return nil, ErrWorkerIDIsRequired
	}
	if c.ClaimedBy == nil || *c.ClaimedBy != workerID {
		return nil, ErrTaskClaimNotOwner
	}

	c.UpdatedAt = time.Now().UTC()
	c.ClaimedBy = nil
	c.ClaimedAt = nil

	return c.validateThenReturn()
}

func (c *TaskClaim) Heartbeat(workerID string) (*TaskClaim, error) {
	if strings.TrimSpace(workerID) == "" {
		return nil, ErrWorkerIDIsRequired
	}
	if c.ClaimedBy == nil || *c.ClaimedBy != workerID {
		return nil, ErrTaskClaimNotOwner
	}

	now := time.Now().UTC()

	c.UpdatedAt = now
	c.ClaimedAt = &now

	return c.validateThenReturn()
}

func (c *TaskClaim) ForceRelease() (*TaskClaim, error) {
	c.UpdatedAt = time.Now().UTC()
	c.ClaimedBy = nil
	c.ClaimedAt = nil

	return c.validateThenReturn()
}

func (c *TaskClaim) IncrementRetryCount() (*TaskClaim, error) {
	c.UpdatedAt = time.Now().UTC()
	c.RetryCount++

	return c.validateThenReturn()
}

func (c *TaskClaim) ResetRetryCount() (*TaskClaim, error) {
	c.UpdatedAt = time.Now().UTC()
	c.RetryCount = 0

	return c.validateThenReturn()
}

func (c *TaskClaim) Delete() (*TaskClaim, error) {
	utcNow := time.Now().UTC()

	c.UpdatedAt = utcNow
	c.DeletedAt = new(utcNow)

	return c.validateThenReturn()
}

func (c *TaskClaim) IsClaimed() bool {
	return c.ClaimedBy != nil && c.ClaimedAt != nil
}

func (c *TaskClaim) IsStale(staleAfter time.Duration) bool {
	if !c.IsClaimed() {
		return false
	}
	return time.Since(*c.ClaimedAt) > staleAfter
}

func (c *TaskClaim) validateThenReturn() (*TaskClaim, error) {
	if c.TaskExecutionID == uuid.Nil {
		return nil, ErrTaskExecutionIDIsRequired
	}
	if c.RetryCount < 0 {
		return nil, ErrRetryCountIsNegative
	}

	return c, nil
}
