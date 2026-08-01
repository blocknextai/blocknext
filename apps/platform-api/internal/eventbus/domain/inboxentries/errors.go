package inboxentries

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrHandlerKeyRequired = apperror.Validation("handler key is required")
)
