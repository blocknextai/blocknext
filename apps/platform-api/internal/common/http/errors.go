package http

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrInvalidRequest   = apperror.Validation("invalid request")
	ErrRateLimitReached = apperror.RateLimited("rate limit exceeded")
)
