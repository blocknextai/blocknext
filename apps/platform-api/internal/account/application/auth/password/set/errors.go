package set

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrPasswordAlreadyExists = apperror.Conflict("password already set; use change-password instead")
	ErrEmailRequired         = apperror.Conflict("email must be linked before setting a password")
)
