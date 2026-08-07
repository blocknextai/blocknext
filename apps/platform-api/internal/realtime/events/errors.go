package events

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrFailedToMarshalTaskEvent = apperror.Internal("failed to marshal task event")
	ErrFailedToMarshalNodeEvent = apperror.Internal("failed to marshal node event")
)
