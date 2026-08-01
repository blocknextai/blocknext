// Package email defines a provider-agnostic interface for sending email
// messages, along with the shared message type and validation used by its
// provider implementations.
package email

import (
	"context"
)

// EmailMessage represents a single email message to be sent, including its
// addressing, subject, body, and delivery options.
type EmailMessage struct {
	From           string
	To             []string
	Cc             []string
	Bcc            []string
	Subject        string
	Body           string
	IsHTML         bool
	SendSeparately bool
}

// EmailSender sends email messages through a concrete provider.
type EmailSender interface {
	// Send delivers the given message, returning an error if delivery fails.
	Send(ctx context.Context, message EmailMessage) error
}
