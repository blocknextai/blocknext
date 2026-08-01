package regenerate

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrTokenRefreshFailed        = apperror.Unavailable("token refresh failed")
	ErrRefreshTokenInvalid       = apperror.Unauthorized("refresh token is invalid; re-authentication required")
	ErrInvalidTokenURL           = apperror.Validation("invalid token url")
	ErrPlatformCredentialMissing = apperror.Internal("platform credential missing")
)
