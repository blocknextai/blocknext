package config

type WebhookOptions struct {
	Trigger WebhookTriggerOptions `envPrefix:"TRIGGER_"`
}

type WebhookTriggerOptions struct {
	URLTemplate string `env:"URL_TEMPLATE"`
}
