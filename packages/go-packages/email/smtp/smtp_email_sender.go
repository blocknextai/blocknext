// Package smtp provides an email.EmailSender implementation that delivers
// messages over an authenticated TLS SMTP connection.
package smtp

import (
	"context"
	"crypto/tls"
	"log/slog"
	"net"
	"net/smtp"
	"strings"

	"github.com/blocknextai/go-packages/cast"
	"github.com/blocknextai/go-packages/email"
)

type smtpEmailSender struct {
	user        string
	password    string
	host        string
	port        int
	defaultFrom string
}

// New returns an email.EmailSender that delivers messages over a TLS SMTP
// connection to host:port authenticated with user and password, applying
// defaultFrom when a message has no From address.
func New(user string, password string, host string, port int, defaultFrom string) email.EmailSender {
	return &smtpEmailSender{
		user:        user,
		password:    password,
		host:        host,
		port:        port,
		defaultFrom: defaultFrom,
	}
}

func (s *smtpEmailSender) Send(ctx context.Context, msg email.EmailMessage) error {
	if strings.TrimSpace(msg.From) == "" {
		msg.From = s.defaultFrom
	}

	if err := email.Validate(msg); err != nil {
		return err
	}

	addr := net.JoinHostPort(s.host, cast.ToString(s.port))

	tlsConfig := &tls.Config{
		ServerName: s.host,
		MinVersion: tls.VersionTLS12,
	}

	dialer := net.Dialer{}
	rawConn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return email.ErrTLSRequired
	}

	conn := tls.Client(rawConn, tlsConfig)
	if err := conn.HandshakeContext(ctx); err != nil {
		if closeErr := rawConn.Close(); closeErr != nil {
			slog.Error("failed to close raw connection",
				"package", "email.smtp",
				"error", closeErr,
			)
		}
		return email.ErrTLSRequired
	}

	if deadline, ok := ctx.Deadline(); ok {
		if err := conn.SetDeadline(deadline); err != nil {
			slog.Warn("failed to set connection deadline",
				"package", "email.smtp",
				"error", err,
			)
		}
	}
	defer func() {
		if err := conn.Close(); err != nil {
			slog.Error("failed to close TLS connection",
				"package", "email.smtp",
				"error", err,
			)
		}
	}()

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return email.ErrSendFailed
	}
	defer func() {
		if err := client.Close(); err != nil {
			slog.Error("failed to close SMTP client",
				"package", "email.smtp",
				"error", err,
			)
		}
	}()

	auth := smtp.PlainAuth("", s.user, s.password, s.host)
	if err := client.Auth(auth); err != nil {
		return email.ErrSendFailed
	}

	toGroups := [][]string{msg.To}
	if msg.SendSeparately {
		toGroups = make([][]string, 0, len(msg.To))
		for _, to := range msg.To {
			toGroups = append(toGroups, []string{to})
		}
	}

	for i, toGroup := range toGroups {
		if i > 0 {
			if err := client.Reset(); err != nil {
				return email.ErrSendFailed
			}
		}

		if err := client.Mail(msg.From); err != nil {
			return email.ErrSendFailed
		}

		for _, to := range toGroup {
			if err := client.Rcpt(to); err != nil {
				return email.ErrSendFailed
			}
		}

		for _, cc := range msg.Cc {
			if err := client.Rcpt(cc); err != nil {
				return email.ErrSendFailed
			}
		}

		for _, bcc := range msg.Bcc {
			if err := client.Rcpt(bcc); err != nil {
				return email.ErrSendFailed
			}
		}

		writer, err := client.Data()
		if err != nil {
			return email.ErrSendFailed
		}

		if _, err := writer.Write(buildMessage(msg, toGroup)); err != nil {
			return email.ErrSendFailed
		}

		if err := writer.Close(); err != nil {
			return email.ErrSendFailed
		}
	}

	return client.Quit()
}

func buildMessage(msg email.EmailMessage, toList []string) []byte {
	var builder strings.Builder

	builder.WriteString("From: ")
	builder.WriteString(msg.From)
	builder.WriteString("\r\n")

	builder.WriteString("To: ")
	builder.WriteString(strings.Join(toList, ", "))
	builder.WriteString("\r\n")

	if len(msg.Cc) > 0 {
		builder.WriteString("Cc: ")
		builder.WriteString(strings.Join(msg.Cc, ", "))
		builder.WriteString("\r\n")
	}

	builder.WriteString("Subject: ")
	builder.WriteString(msg.Subject)
	builder.WriteString("\r\n")

	builder.WriteString("MIME-Version: 1.0\r\n")

	if msg.IsHTML {
		builder.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	} else {
		builder.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	}

	builder.WriteString("\r\n")
	builder.WriteString(msg.Body)

	return []byte(builder.String())
}
