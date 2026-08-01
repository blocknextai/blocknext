// Package apperror defines a structured application error type with a
// classification kind that maps to HTTP status codes.
package apperror

import (
	"errors"
)

// Kind classifies an application error and determines its HTTP status mapping.
type Kind int

// The set of application error kinds.
const (
	KindUnknown Kind = iota
	KindValidation
	KindNotFound
	KindConflict
	KindUnauthorized
	KindForbidden
	KindPayloadTooLarge
	KindUnsupportedMedia
	KindRateLimited
	KindInternal
	KindBadGateway
	KindUnavailable
	KindGatewayTimeout
)

// Error is a structured application error carrying a kind, optional code,
// message and underlying cause.
type Error struct {
	Kind    Kind   `json:"kind"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"cause"`
}

// Error returns the error message, or an empty string for a nil receiver.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// Unwrap returns the underlying cause, enabling errors.Is and errors.As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// Is reports whether target is an Error with the same kind, code and message.
func (e *Error) Is(target error) bool {
	if e == nil || target == nil {
		return e == nil && target == nil
	}
	other, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Kind == other.Kind && e.Code == other.Code && e.Message == other.Message
}

// WithCause returns a copy of the error with the given cause set.
func (e *Error) WithCause(cause error) *Error {
	return &Error{Kind: e.Kind, Code: e.Code, Message: e.Message, Cause: cause}
}

// WithCode returns a copy of the error with the given code set.
func (e *Error) WithCode(code string) *Error {
	return &Error{Kind: e.Kind, Code: code, Message: e.Message, Cause: e.Cause}
}

// Validation returns an Error of kind KindValidation with the given message.
func Validation(message string) *Error {
	return &Error{Kind: KindValidation, Message: message}
}

// NotFound returns an Error of kind KindNotFound with the given message.
func NotFound(message string) *Error {
	return &Error{Kind: KindNotFound, Message: message}
}

// Conflict returns an Error of kind KindConflict with the given message.
func Conflict(message string) *Error {
	return &Error{Kind: KindConflict, Message: message}
}

// Unauthorized returns an Error of kind KindUnauthorized with the given message.
func Unauthorized(message string) *Error {
	return &Error{Kind: KindUnauthorized, Message: message}
}

// Forbidden returns an Error of kind KindForbidden with the given message.
func Forbidden(message string) *Error {
	return &Error{Kind: KindForbidden, Message: message}
}

// PayloadTooLarge returns an Error of kind KindPayloadTooLarge with the given
// message.
func PayloadTooLarge(message string) *Error {
	return &Error{Kind: KindPayloadTooLarge, Message: message}
}

// UnsupportedMedia returns an Error of kind KindUnsupportedMedia with the given
// message.
func UnsupportedMedia(message string) *Error {
	return &Error{Kind: KindUnsupportedMedia, Message: message}
}

// RateLimited returns an Error of kind KindRateLimited with the given message.
func RateLimited(message string) *Error {
	return &Error{Kind: KindRateLimited, Message: message}
}

// Internal returns an Error of kind KindInternal with the given message.
func Internal(message string) *Error {
	return &Error{Kind: KindInternal, Message: message}
}

// BadGateway returns an Error of kind KindBadGateway with the given message.
func BadGateway(message string) *Error {
	return &Error{Kind: KindBadGateway, Message: message}
}

// Unavailable returns an Error of kind KindUnavailable with the given message.
func Unavailable(message string) *Error {
	return &Error{Kind: KindUnavailable, Message: message}
}

// GatewayTimeout returns an Error of kind KindGatewayTimeout with the given
// message.
func GatewayTimeout(message string) *Error {
	return &Error{Kind: KindGatewayTimeout, Message: message}
}

// As reports whether err is an *Error and returns it if so.
func As(err error) (*Error, bool) {
	if appErr, ok := errors.AsType[*Error](err); ok {
		return appErr, true
	}
	return nil, false
}
