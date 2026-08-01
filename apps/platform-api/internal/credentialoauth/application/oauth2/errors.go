package oauth2

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrInvalidState                         = apperror.Validation("invalid state")
	ErrInvalidCredential                    = apperror.Validation("invalid credential")
	ErrInvalidCredentialData                = apperror.Validation("invalid credential data")
	ErrTokenExchangeFailed                  = apperror.Internal("token exchange failed")
	ErrPlatformCredentialMissing            = apperror.Internal("platform credential missing")
	ErrPlatformCredentialMissingClientCreds = apperror.Internal("platform credential missing client credentials")
)
