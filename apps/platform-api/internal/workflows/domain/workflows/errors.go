package workflows

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrWorkflowNotFound         = apperror.NotFound("workflow not found")
	ErrOrganizationIDIsRequired = apperror.Validation("organization id is required")
	ErrOwnerIDIsRequired        = apperror.Validation("owner id is required")
	ErrTitleIsRequired          = apperror.Validation("title is required")
	ErrTitleTooLong             = apperror.Validation("title too long")
)
