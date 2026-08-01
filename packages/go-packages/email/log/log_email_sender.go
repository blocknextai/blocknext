// Package log provides an email.EmailSender implementation that logs messages
// via slog instead of delivering them, useful for development and testing.
package log

import (
	"context"
	"log/slog"
	"strings"

	"github.com/blocknextai/go-packages/email"
)

type logEmailSender struct {
	logger      *slog.Logger
	defaultFrom string
}

// New returns an email.EmailSender that logs each message via the given
// logger, falling back to slog.Default when logger is nil and using
// defaultFrom when a message has no From address.
func New(logger *slog.Logger, defaultFrom string) email.EmailSender {
	if logger == nil {
		logger = slog.Default()
	}

	return &logEmailSender{
		logger:      logger,
		defaultFrom: defaultFrom,
	}
}

func (s *logEmailSender) Send(ctx context.Context, msg email.EmailMessage) error {
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
		s.logger.LogAttrs(ctx, slog.LevelInfo, "email sent",
			slog.String("provider", "log"),
			slog.String("from", msg.From),
			slog.Any("to", toGroup),
			slog.Any("cc", msg.Cc),
			slog.Any("bcc", msg.Bcc),
			slog.String("subject", msg.Subject),
			slog.String("body", msg.Body),
			slog.Bool("is_html", msg.IsHTML),
		)
	}

	return nil
}
