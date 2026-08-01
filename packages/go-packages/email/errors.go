package email

import (
	"errors"
)

var (
	// ErrInvalidFromAddress indicates the From address is missing or malformed.
	ErrInvalidFromAddress = errors.New("invalid from address")
	// ErrInvalidToAddress indicates a To address is missing or malformed.
	ErrInvalidToAddress = errors.New("invalid to address")
	// ErrInvalidCcAddress indicates a Cc address is malformed.
	ErrInvalidCcAddress = errors.New("invalid cc address")
	// ErrInvalidBccAddress indicates a Bcc address is malformed.
	ErrInvalidBccAddress = errors.New("invalid bcc address")
	// ErrEmptySubject indicates the message subject is empty.
	ErrEmptySubject = errors.New("empty subject")
	// ErrEmptyBody indicates the message body is empty.
	ErrEmptyBody = errors.New("empty body")
	// ErrTLSRequired indicates a secure TLS connection could not be established.
	ErrTLSRequired = errors.New("tls required")
	// ErrSendFailed indicates the message could not be delivered by the provider.
	ErrSendFailed = errors.New("send failed")
)
