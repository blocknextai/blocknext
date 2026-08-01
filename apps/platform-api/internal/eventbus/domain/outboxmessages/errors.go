package outboxmessages

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrEventNameRequired       = apperror.Validation("event name is required")
	ErrPayloadRequired         = apperror.Validation("payload is required")
	ErrInvalidStatusTransition = apperror.Internal("invalid outbox message status transition")
)
