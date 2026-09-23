package adapters

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrAdapterNotFound  = apperror.NotFound("adapter not found")
	ErrInvalidSignature = apperror.Unauthorized("invalid signature")
)
