package refreshtokens

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrRefreshTokenNotFound         = apperror.NotFound("refresh token not found")
	ErrRefreshTokenUsed             = apperror.Validation("refresh token already used")
	ErrRefreshTokenRevoked          = apperror.Validation("refresh token revoked")
	ErrRefreshTokenExpired          = apperror.Validation("refresh token expired")
	ErrRefreshTokenClientMismatch   = apperror.Validation("refresh token was issued to another client")
	ErrRefreshTokenResourceMismatch = apperror.Validation("refresh token was issued for another resource")
	ErrMissingRefreshToken          = apperror.Validation("refresh token is required")
	ErrInvalidGrantID               = apperror.Validation("invalid grant id")
	ErrInvalidClientID              = apperror.Validation("invalid client id")
	ErrInvalidTokenHash             = apperror.Validation("invalid token hash")
	ErrInvalidResource              = apperror.Validation("invalid resource")
	ErrInvalidScopes                = apperror.Validation("invalid scopes")
)
