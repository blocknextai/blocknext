package tokens

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrInvalidAccessToken         = apperror.Unauthorized("invalid access token")
	ErrExpiredAccessToken         = apperror.Unauthorized("expired access token")
	ErrInvalidAccessTokenAudience = apperror.Unauthorized("access token was not issued for this resource")
	ErrInvalidSigningKey          = apperror.Internal("invalid access token signing key")
	ErrInvalidIssuer              = apperror.Internal("invalid access token issuer")
	ErrFailedToSignAccessToken    = apperror.Internal("failed to sign access token")
	ErrFailedToGenerateToken      = apperror.Internal("failed to generate token")
	ErrFailedToExchangeToken      = apperror.Internal("failed to exchange token")
)
