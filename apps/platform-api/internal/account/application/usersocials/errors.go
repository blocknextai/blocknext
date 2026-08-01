package usersocials

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrURLRequired   = apperror.Validation("url is required")
	ErrHttpsRequired = apperror.Validation("https is required")
	ErrInvalidURL    = apperror.Validation("invalid url")
)
