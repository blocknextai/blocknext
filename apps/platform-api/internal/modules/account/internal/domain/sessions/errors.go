package sessions

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrSessionNotFound       = apperror.NotFound("session not found")
	ErrRefreshTokenExpired   = apperror.Unauthorized("refresh token expired")
	ErrRefreshTokenReused    = apperror.Unauthorized("refresh token reused")
	ErrSessionAlreadyRevoked = apperror.Conflict("session already revoked")
)
