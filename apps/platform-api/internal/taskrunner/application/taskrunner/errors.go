package taskrunner

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrNodeExecutionNotFound   = apperror.NotFound("node execution not found")
	ErrTaskAlreadyCompleted    = apperror.Conflict("task already completed")
	ErrNodeExecutionIDNotFound = apperror.NotFound("node execution id not found")
	ErrExecutorNotFound        = apperror.NotFound("executor not found")
	ErrTaskNotFound            = apperror.NotFound("task not found")
	ErrTaskNotFailed           = apperror.Conflict("task not failed")
	ErrTaskDAGIsEmpty          = apperror.Validation("task dag is empty")
	ErrNoStartNodesFound       = apperror.Validation("no start nodes found")
	ErrNoNodesFoundInWorkflow  = apperror.Validation("no nodes found in workflow")
	ErrNoNodeExecutionsFound   = apperror.NotFound("no node executions found")
	ErrSomeNodesFailed         = apperror.Internal("some nodes failed")
	ErrNodeDisabled            = apperror.Conflict("node disabled")
	ErrTaskCancelled           = apperror.Internal("task cancelled")
	ErrTaskExecutionPanic      = apperror.Internal("task execution panic")
	ErrWorkerPoolFull          = apperror.Unavailable("worker pool full")
	ErrTaskClaimNotAvailable   = apperror.Conflict("task claim not available")
)
