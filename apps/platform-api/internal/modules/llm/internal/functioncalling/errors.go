package functioncalling

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrEmptyFunctionDeclarations = apperror.Validation("empty function declarations")
	ErrEmptyData                 = apperror.Validation("empty data")
	ErrNoCandidates              = apperror.Internal("no candidates")
	ErrNoContentParts            = apperror.Internal("no content parts")
	ErrNoFunctionCallName        = apperror.Internal("no function call name")
	ErrRateLimited               = apperror.RateLimited("rate limited")
	ErrProviderRequestFailed     = apperror.Unavailable("provider request failed")
)
