package domain

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrInvalidRequest = apperror.Validation("invalid request")
)
