package publishing

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrMarshalEvent = apperror.Internal("eventbus: failed to marshal event")
)
