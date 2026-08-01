package validation

import (
	"net/mail"
	"strings"

	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrEmailIsRequired = apperror.Validation("email is required")
	ErrInvalidEmail    = apperror.Validation("invalid email")
)

func Email(email string) error {
	if strings.TrimSpace(email) == "" {
		return ErrEmailIsRequired
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidEmail
	}
	return nil
}
