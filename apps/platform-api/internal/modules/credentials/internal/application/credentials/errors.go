package credentials

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrFailedToCreateCredential       = apperror.Internal("failed to create credential")
	ErrFailedToUpdateCredential       = apperror.Internal("failed to update credential")
	ErrFailedToDeleteCredential       = apperror.Internal("failed to delete credential")
	ErrFailedToGetCredentialByID      = apperror.Internal("failed to get credential by id")
	ErrFailedToEncryptData            = apperror.Internal("failed to encrypt data")
	ErrFailedToDecryptData            = apperror.Internal("failed to decrypt data")
	ErrPlatformCredentialNotSupported = apperror.Validation("platform credential not supported")
)
