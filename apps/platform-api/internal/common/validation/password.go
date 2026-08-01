package validation

import (
	"github.com/blocknextai/go-packages/apperror"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 128
)

var (
	ErrPasswordIsRequired = apperror.Validation("password is required")
	ErrPasswordTooShort   = apperror.Validation("password must be at least 8 characters")
	ErrPasswordTooLong    = apperror.Validation("password must be at most 128 characters")
)

func Password(s string) error {
	if s == "" {
		return ErrPasswordIsRequired
	}
	return nil
}

func NewPassword(s string) error {
	if err := Password(s); err != nil {
		return err
	}
	if len(s) < MinPasswordLength {
		return ErrPasswordTooShort
	}
	if len(s) > MaxPasswordLength {
		return ErrPasswordTooLong
	}
	return nil
}
