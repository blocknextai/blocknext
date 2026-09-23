package linkedaccounts

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrFailedToGetUser                  = apperror.Internal("failed to get user")
	ErrFailedToAddLinkedAccount         = apperror.Internal("failed to add linked account")
	ErrFailedToDeleteLinkedAccount      = apperror.Internal("failed to delete linked account")
	ErrLinkedAccountAlreadyExists       = apperror.Conflict("linked account already exists")
	ErrCannotDeletePrimaryLinkedAccount = apperror.Forbidden("cannot delete primary linked account")
)
