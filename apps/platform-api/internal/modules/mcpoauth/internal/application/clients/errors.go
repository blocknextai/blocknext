package clients

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrFailedToGenerateClientCredentials = apperror.Internal("failed to generate client credentials")
	ErrFailedToRegisterClient            = apperror.Internal("failed to register client")
)
