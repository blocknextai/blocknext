package email

import (
	"net/mail"
	"strings"
)

// Validate checks that the message has a valid From address, at least one
// valid To address, valid Cc and Bcc addresses, and a non-empty subject and
// body, returning the corresponding error for the first violation found.
func Validate(msg EmailMessage) error {
	if err := validateAddress(msg.From); err != nil {
		return ErrInvalidFromAddress
	}

	if len(msg.To) == 0 {
		return ErrInvalidToAddress
	}

	for _, to := range msg.To {
		if err := validateAddress(to); err != nil {
			return ErrInvalidToAddress
		}
	}

	for _, cc := range msg.Cc {
		if err := validateAddress(cc); err != nil {
			return ErrInvalidCcAddress
		}
	}

	for _, bcc := range msg.Bcc {
		if err := validateAddress(bcc); err != nil {
			return ErrInvalidBccAddress
		}
	}

	if strings.TrimSpace(msg.Subject) == "" {
		return ErrEmptySubject
	}

	if strings.TrimSpace(msg.Body) == "" {
		return ErrEmptyBody
	}

	return nil
}

func validateAddress(addr string) error {
	_, err := mail.ParseAddress(addr)
	return err
}
