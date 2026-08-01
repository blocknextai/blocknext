package nodeexecutions

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrNodeExecutionNotFound = apperror.NotFound("node execution not found")
	ErrTaskIDIsRequired      = apperror.Validation("task id is required")
	ErrNodeTypeIsRequired    = apperror.Validation("node type is required")
	ErrNodeIDIsRequired      = apperror.Validation("node id is required")
	ErrStatusIsRequired      = apperror.Validation("status is required")
)
