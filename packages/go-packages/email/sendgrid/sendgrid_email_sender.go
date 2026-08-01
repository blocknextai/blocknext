// Package sendgrid provides an email.EmailSender implementation backed by the
// SendGrid HTTP API.
package sendgrid

import (
	"context"
	"strings"

	"github.com/blocknextai/go-packages/email"
	"github.com/blocknextai/go-packages/httpclient"
)

type sendgridEmailSender struct {
	client      *httpclient.Client
	defaultFrom string
}

// New returns an email.EmailSender that delivers messages through the SendGrid
// API using the given API key, applying defaultFrom when a message has no
// From address.
func New(apiKey string, defaultFrom string) email.EmailSender {
	client := httpclient.NewClientBuilder().
		BaseURL("https://api.sendgrid.com").
		BearerAuth(apiKey).
		JSONContentType().
		Build()

	return &sendgridEmailSender{
		client:      client,
		defaultFrom: defaultFrom,
	}
}

func (s *sendgridEmailSender) Send(ctx context.Context, msg email.EmailMessage) error {
	if strings.TrimSpace(msg.From) == "" {
		msg.From = s.defaultFrom
	}

	if err := email.Validate(msg); err != nil {
		return err
	}

	contentType := "text/plain"
	if msg.IsHTML {
		contentType = "text/html"
	}

	var personalizations []sendgridPersonalization
	if msg.SendSeparately {
		personalizations = make([]sendgridPersonalization, 0, len(msg.To))
		for _, to := range msg.To {
			p := sendgridPersonalization{
				To: []sendgridAddress{{Email: to}},
			}
			if len(msg.Cc) > 0 {
				p.Cc = toAddressList(msg.Cc)
			}
			if len(msg.Bcc) > 0 {
				p.Bcc = toAddressList(msg.Bcc)
			}
			personalizations = append(personalizations, p)
		}
	} else {
		p := sendgridPersonalization{
			To: toAddressList(msg.To),
		}
		if len(msg.Cc) > 0 {
			p.Cc = toAddressList(msg.Cc)
		}
		if len(msg.Bcc) > 0 {
			p.Bcc = toAddressList(msg.Bcc)
		}
		personalizations = []sendgridPersonalization{p}
	}

	body := sendgridRequest{
		Personalizations: personalizations,
		From:             sendgridAddress{Email: msg.From},
		Subject:          msg.Subject,
		Content: []sendgridContent{
			{Type: contentType, Value: msg.Body},
		},
	}

	errorResponse := &sendgridErrorResponse{}

	response, err := s.client.Post("/v3/mail/send").
		Context(ctx).
		Body(body).
		Do(nil, errorResponse)
	if err != nil {
		return email.ErrSendFailed
	}

	if !response.IsSuccess() {
		return email.ErrSendFailed
	}

	return nil
}

func toAddressList(addresses []string) []sendgridAddress {
	result := make([]sendgridAddress, 0, len(addresses))
	for _, addr := range addresses {
		result = append(result, sendgridAddress{Email: addr})
	}
	return result
}
