package credentials

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrCredentialNotFound = apperror.NotFound("credential not found")
	ErrOwnerTypeIsInvalid = apperror.Validation("owner type is invalid")
	ErrOwnerIDIsRequired  = apperror.Validation("owner id is required")
	ErrTitleIsRequired    = apperror.Validation("title is required")
	ErrKeyIsRequired      = apperror.Validation("key is required")
	ErrDataIsRequired     = apperror.Validation("data is required")
	ErrInvalidSourceType  = apperror.Validation("invalid source type")
	ErrInvalidKey         = apperror.Validation("invalid key")
	ErrInvalidTitle       = apperror.Validation("invalid title")
	ErrInvalidData        = apperror.Validation("invalid data")
	ErrKeyTooLong         = apperror.Validation("key too long")
	ErrTitleTooLong       = apperror.Validation("title too long")
)
