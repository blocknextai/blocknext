package domain

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrInvalidURL      = apperror.Validation("invalid url")
	ErrInvalidScheme   = apperror.Validation("invalid scheme")
	ErrBlockedIP       = apperror.Forbidden("blocked ip")
	ErrFileTooLarge    = apperror.PayloadTooLarge("file too large")
	ErrDownloadFailed  = apperror.BadGateway("download failed")
	ErrDownloadTimeout = apperror.GatewayTimeout("download timeout")
)
