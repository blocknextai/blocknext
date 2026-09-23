package workflows

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrFailedToUpdateWorkflow = apperror.Internal("failed to update workflow")
	ErrFailedToDeleteWorkflow = apperror.Internal("failed to delete workflow")
)
