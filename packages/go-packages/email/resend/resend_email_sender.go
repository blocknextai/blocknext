// Package resend provides an email.EmailSender implementation backed by the
// Resend HTTP API.
package resend

import (
	"context"
	"strings"

	"github.com/blocknextai/go-packages/email"
	"github.com/blocknextai/go-packages/httpclient"
)

type resendEmailSender struct {
	client      *httpclient.Client
	defaultFrom string
}

// New returns an email.EmailSender that delivers messages through the Resend
// API using the given API key, applying defaultFrom when a message has no
// From address.
func New(apiKey string, defaultFrom string) email.EmailSender {
	client := httpclient.NewClientBuilder().
		BaseURL("https://api.resend.com").
		BearerAuth(apiKey).
		JSONContentType().
		Build()

	return &resendEmailSender{
		client:      client,
		defaultFrom: defaultFrom,
	}
}

func (s *resendEmailSender) Send(ctx context.Context, msg email.EmailMessage) error {
	if strings.TrimSpace(msg.From) == "" {
		msg.From = s.defaultFrom
	}

	if err := email.Validate(msg); err != nil {
		return err
	}

	toGroups := [][]string{msg.To}
	if msg.SendSeparately {
		toGroups = make([][]string, 0, len(msg.To))
		for _, to := range msg.To {
			toGroups = append(toGroups, []string{to})
		}
	}

	for _, toGroup := range toGroups {
		body := resendRequest{
			From:    msg.From,
			To:      toGroup,
			Cc:      msg.Cc,
			Bcc:     msg.Bcc,
			Subject: msg.Subject,
		}

		if msg.IsHTML {
			body.HTML = msg.Body
		} else {
			body.Text = msg.Body
		}

		errorResponse := &resendErrorResponse{}

		response, err := s.client.Post("/emails").
			Context(ctx).
			Body(body).
			Do(nil, errorResponse)
		if err != nil {
			return email.ErrSendFailed
		}

		if !response.IsSuccess() {
			return email.ErrSendFailed
		}
	}

	return nil
}
