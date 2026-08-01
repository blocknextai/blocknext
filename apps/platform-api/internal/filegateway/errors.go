package filegateway

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrDownloadFailed     = apperror.Internal("file download failed")
	ErrInvalidDownloadURL = apperror.Validation("download url is not allowed")
)
