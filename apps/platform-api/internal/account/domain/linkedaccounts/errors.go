package linkedaccounts

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrLinkedAccountNotFound          = apperror.NotFound("linked account not found")
	ErrUserIDIsRequired               = apperror.Validation("user id is required")
	ErrProviderIDIsRequired           = apperror.Validation("provider id is required")
	ErrIdentifierIsRequired           = apperror.Validation("identifier is required")
	ErrLinkedAccountProviderTaken     = apperror.Conflict("linked account provider already taken")
	ErrLinkedAccountUserProviderTaken = apperror.Conflict("linked account already exists for user and provider")
)
