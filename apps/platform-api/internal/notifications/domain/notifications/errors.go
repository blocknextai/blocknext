package notifications

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrInvalidType         = apperror.Validation("invalid notification type")
	ErrInvalidLevel        = apperror.Validation("invalid notification level")
	ErrInvalidAudienceType = apperror.Validation("invalid notification audience type")
	ErrInvalidAudienceID   = apperror.Validation("invalid notification audience id")
	ErrInvalidTitle        = apperror.Validation("invalid notification title")
)
