package domain

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrInvalidUploadID    = apperror.Validation("invalid upload id")
	ErrUploadRuleNotFound = apperror.NotFound("upload rule not found")
	ErrInvalidFile        = apperror.Validation("invalid file")
	ErrInvalidUpload      = apperror.Validation("invalid upload")
	ErrInvalidFilename    = apperror.Validation("invalid filename")
	ErrFileNotRead        = apperror.Internal("file not read")
	ErrStorageError       = apperror.Internal("storage error")
	ErrMaxSizeExceeded    = apperror.PayloadTooLarge("max size exceeded")
	ErrMIMENotAllowed     = apperror.UnsupportedMedia("mime type not allowed")
)
