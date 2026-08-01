package sessions

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrSessionNotFound = apperror.NotFound("session not found")
	ErrTitleIsRequired = apperror.Validation("title is required")
)
