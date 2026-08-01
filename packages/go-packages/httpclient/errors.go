package httpclient

import (
	"errors"
)

var (
	// ErrInvalidBody indicates that the request body does not match the configured content type.
	ErrInvalidBody = errors.New("invalid body")
	// ErrMaxRetriesExceeded indicates that all retry attempts were exhausted without success.
	ErrMaxRetriesExceeded = errors.New("max retries exceeded")
	// ErrInvalidURL indicates that the resolved request URL could not be parsed.
	ErrInvalidURL = errors.New("invalid url")
)
