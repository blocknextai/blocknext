package application

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrFailedToEncryptSecret = apperror.Internal("failed to encrypt secret")
)
