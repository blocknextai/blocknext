package taskclaims

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrTaskExecutionIDIsRequired = apperror.Validation("task execution id is required")
	ErrWorkerIDIsRequired        = apperror.Validation("worker id is required")
	ErrRetryCountIsNegative      = apperror.Validation("retry count cannot be negative")
	ErrTaskClaimNotFound         = apperror.NotFound("task claim not found")
	ErrTaskClaimAlreadyClaimed   = apperror.Conflict("task claim already claimed")
	ErrTaskClaimNotOwner         = apperror.Conflict("task claim not owned by worker")
)
