package config

type EmailProvider string

const (
	EmailProviderLog      EmailProvider = "log"
	EmailProviderSMTP     EmailProvider = "smtp"
	EmailProviderResend   EmailProvider = "resend"
	EmailProviderSendGrid EmailProvider = "sendgrid"
)

type EmailSenderOptions struct {
	Enabled  bool                       `env:"ENABLED"`
	Provider EmailProvider              `env:"PROVIDER"`
	From     string                     `env:"FROM"`
	SMTP     EmailSenderSMTPOptions     `envPrefix:"SMTP_"`
	Resend   EmailSenderResendOptions   `envPrefix:"RESEND_"`
	SendGrid EmailSenderSendGridOptions `envPrefix:"SENDGRID_"`
}

type EmailSenderSMTPOptions struct {
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
}

type EmailSenderResendOptions struct {
	APIKey string `env:"API_KEY"`
}

type EmailSenderSendGridOptions struct {
	APIKey string `env:"API_KEY"`
}
