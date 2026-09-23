package organizationusers

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrOrganizationUserNotFound      = apperror.NotFound("organization user not found")
	ErrOrganizationIDIsRequired      = apperror.Validation("organization id is required")
	ErrUserIDIsRequired              = apperror.Validation("user id is required")
	ErrRoleIsRequired                = apperror.Validation("role is required")
	ErrInvalidRole                   = apperror.Validation("invalid role")
	ErrOrganizationUserAlreadyExists = apperror.Conflict("organization user already exists")
)
