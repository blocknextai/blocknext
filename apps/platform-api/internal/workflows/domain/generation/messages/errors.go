package messages

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrSessionIDIsRequired = apperror.Validation("session id is required")
)
