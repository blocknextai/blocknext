package mailer

import (
	"context"
	"strings"

	pkgEmail "github.com/blocknextai/go-packages/email"
)

const (
	verifyEmailPath        = "/auth/verify-email"
	resetPasswordPath      = "/auth/reset-password"
	emailChangeConfirmPath = "/auth/email-change/confirm"
	magicLinkPath          = "/auth/magic-link/callback"
	loginPath              = "/auth/login"
)

type Mailer struct {
	sender            pkgEmail.EmailSender
	platformUIBaseURL string
}

func NewMailer(sender pkgEmail.EmailSender, platformUIBaseURL string) *Mailer {
	return &Mailer{
		sender:            sender,
		platformUIBaseURL: platformUIBaseURL,
	}
}

func (m *Mailer) SendVerification(ctx context.Context, to, plainToken string) error {
	body := "Welcome! Please verify your email address by clicking the link below:\n\n" +
		m.buildLink(verifyEmailPath, plainToken) +
		"\n\nThis link will expire in 24 hours."

	return m.send(ctx, to, "Verify your email address", body)
}

func (m *Mailer) SendPasswordReset(ctx context.Context, to, plainToken string) error {
	body := "We received a request to reset your password. Click the link below to choose a new one:\n\n" +
		m.buildLink(resetPasswordPath, plainToken) +
		"\n\nIf you did not request this, you can safely ignore this email. The link expires in 1 hour."

	return m.send(ctx, to, "Reset your password", body)
}

func (m *Mailer) SendEmailChangeConfirm(ctx context.Context, to, plainToken string) error {
	body := "Please confirm your new email address by clicking the link below:\n\n" +
		m.buildLink(emailChangeConfirmPath, plainToken) +
		"\n\nIf you did not request this change, you can safely ignore this email. The link expires in 1 hour."

	return m.send(ctx, to, "Confirm your new email address", body)
}

func (m *Mailer) SendMagicLink(ctx context.Context, to, plainToken string) error {
	body := "Click the link below to sign in:\n\n" +
		m.buildLink(magicLinkPath, plainToken) +
		"\n\nThis link will expire in 15 minutes and can only be used once.\nIf you did not request this, you can safely ignore this email."

	return m.send(ctx, to, "Your sign-in link", body)
}

func (m *Mailer) SendAlreadyRegistered(ctx context.Context, to string) error {
	var link strings.Builder
	link.WriteString(m.platformUIBaseURL)
	link.WriteString(loginPath)

	body := "Someone just tried to register a new account with this email address.\n\n" +
		"If this was you, your account already exists. You can sign in here:\n\n" +
		link.String() +
		"\n\nIf it wasn't you, you can safely ignore this email."

	return m.send(ctx, to, "Account already exists", body)
}

func (m *Mailer) SendWelcome(ctx context.Context, to string) error {
	var link strings.Builder
	link.WriteString(m.platformUIBaseURL)
	link.WriteString(loginPath)

	body := "Welcome! Your account is ready.\n\n" +
		"Sign in to start building your first workflow:\n\n" +
		link.String()

	return m.send(ctx, to, "Welcome", body)
}

func (m *Mailer) buildLink(path, plainToken string) string {
	var link strings.Builder
	link.WriteString(m.platformUIBaseURL)
	link.WriteString(path)
	link.WriteString("?token=")
	link.WriteString(plainToken)
	return link.String()
}

func (m *Mailer) send(ctx context.Context, to, subject, body string) error {
	if m.sender == nil {
		return nil
	}

	return m.sender.Send(ctx, pkgEmail.EmailMessage{
		To:      []string{to},
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	})
}
