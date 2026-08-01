package webhooks

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrTriggerNotFound       = apperror.NotFound("trigger not found")
	ErrTriggerInactive       = apperror.Conflict("trigger inactive")
	ErrFailedToDecryptSecret = apperror.Internal("failed to decrypt secret")
)
