package auth

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrFailedToGetUser                   = apperror.Internal("failed to get user")
	ErrFailedToGenerateToken             = apperror.Internal("failed to generate token")
	ErrFailedToHashPassword              = apperror.Internal("failed to hash password")
	ErrFailedToGenerateVerificationToken = apperror.Internal("failed to generate verification token")
	ErrInvalidRefreshToken               = apperror.Unauthorized("invalid refresh token")
	ErrInvalidEmailOrPassword            = apperror.Unauthorized("invalid email or password")
	ErrEmailNotVerified                  = apperror.Forbidden("email not verified")
	ErrPasswordAuthDisabled              = apperror.Forbidden("password authentication is disabled")
	ErrPasswordNotSet                    = apperror.Forbidden("password not set")
	ErrUserBanned                        = apperror.Forbidden("user banned")
)
