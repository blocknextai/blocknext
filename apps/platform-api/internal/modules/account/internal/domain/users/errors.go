package users

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrUserNotFound      = apperror.NotFound("user not found")
	ErrRoleIsRequired    = apperror.Validation("role is required")
	ErrInvalidRole       = apperror.Validation("invalid role")
	ErrEmailAlreadyInUse = apperror.Conflict("email already in use")
)
