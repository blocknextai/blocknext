package passwordcredentials

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrPasswordCredentialNotFound = apperror.NotFound("password credential not found")
	ErrPasswordHashIsRequired     = apperror.Validation("password hash is required")
)
