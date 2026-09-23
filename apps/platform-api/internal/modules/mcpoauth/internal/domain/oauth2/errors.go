package oauth2

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrInvalidBaseURL       = apperror.Validation("issuer and resource urls must be absolute and without a path")
	ErrMalformedRedirectURI = apperror.Validation("malformed redirect uri")
	ErrUnsupportedGrantType = apperror.Validation("unsupported grant type")
)
