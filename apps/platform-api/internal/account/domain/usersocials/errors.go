package usersocials

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrSocialNotFound     = apperror.NotFound("social not found")
	ErrURLIsRequired      = apperror.Validation("url is required")
	ErrUserIDIsRequired   = apperror.Validation("user id is required")
	ErrPlatformIsRequired = apperror.Validation("platform is required")
)
