package toolinvocations

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrToolInvocationNotFound = apperror.NotFound("tool invocation not found")
	ErrOrganizationIDRequired = apperror.Validation("organization id is required")
	ErrSourceInvalid          = apperror.Validation("invalid tool invocation source")
	ErrToolIDIsRequired       = apperror.Validation("tool id is required")
	ErrStatusInvalid          = apperror.Validation("invalid tool invocation status")
)
