package getallwebhooksources

type WebhookSourceResponse struct {
	Source               string `json:"source"`
	Name                 string `json:"name"`
	SupportsVerification bool   `json:"supportsVerification"`
	WebhookURL           string `json:"webhookUrl"`
}

type GetAllWebhookSourcesResponse = []WebhookSourceResponse
