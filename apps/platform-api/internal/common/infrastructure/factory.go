package infrastructure

import (
	"log/slog"

	pkgEmail "github.com/blocknextai/go-packages/email"
	emailLog "github.com/blocknextai/go-packages/email/log"
	emailResend "github.com/blocknextai/go-packages/email/resend"
	emailSendGrid "github.com/blocknextai/go-packages/email/sendgrid"
	emailSmtp "github.com/blocknextai/go-packages/email/smtp"
	"github.com/blocknextai/go-packages/hashing"
	hashingBcrypt "github.com/blocknextai/go-packages/hashing/bcrypt"
	"github.com/blocknextai/platform-api/internal/config"
)

const (
	defaultBcryptCost = 12
)

func NewPasswordHasher(cost int) hashing.Hasher {
	if cost <= 0 {
		cost = defaultBcryptCost
	}
	return hashingBcrypt.New(cost)
}

func NewEmailSender(options config.EmailSenderOptions) pkgEmail.EmailSender {
	if !options.Enabled {
		slog.Info("email sender is disabled")
		return nil
	}

	if options.Provider == "" || options.Provider == config.EmailProviderLog {
		slog.Info("email sender initialized", "provider", config.EmailProviderLog)
		return emailLog.New(slog.Default(), options.From)
	}

	if options.Provider == config.EmailProviderSMTP {
		slog.Info("email sender initialized", "provider", options.Provider)
		return emailSmtp.New(options.SMTP.User, options.SMTP.Password, options.SMTP.Host, options.SMTP.Port, options.From)
	}

	if options.Provider == config.EmailProviderResend {
		slog.Info("email sender initialized", "provider", options.Provider)
		return emailResend.New(options.Resend.APIKey, options.From)
	}

	if options.Provider == config.EmailProviderSendGrid {
		slog.Info("email sender initialized", "provider", options.Provider)
		return emailSendGrid.New(options.SendGrid.APIKey, options.From)
	}

	slog.Warn("unknown email provider — falling back to log",
		"component", "EmailSenderFactory",
		"provider", options.Provider)
	return emailLog.New(slog.Default(), options.From)
}
