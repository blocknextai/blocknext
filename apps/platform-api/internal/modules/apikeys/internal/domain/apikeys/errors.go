package apikeys

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrAPIKeyNotFound      = apperror.NotFound("api key not found")
	ErrInvalidAPIKeyName   = apperror.Validation("invalid api key name")
	ErrInvalidAPIKeyScopes = apperror.Validation("invalid api key scopes")
	ErrInvalidOwnerType    = apperror.Validation("invalid owner type")
	ErrInvalidOwnerID      = apperror.Validation("invalid owner id")
	ErrInvalidAPIKey       = apperror.Validation("invalid api key")
	ErrAPIKeyAlreadyExists = apperror.Conflict("api key already exists")
)
