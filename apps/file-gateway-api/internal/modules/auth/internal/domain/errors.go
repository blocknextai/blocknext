package domain

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrInvalidServiceKey = apperror.Unauthorized("invalid service key")
	ErrInvalidToken      = apperror.Unauthorized("invalid token")
	ErrMissingToken      = apperror.Unauthorized("missing token")
)
