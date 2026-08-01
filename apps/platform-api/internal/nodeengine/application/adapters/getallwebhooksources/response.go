package getallwebhooksources

type WebhookSourceResponse struct {
	Source     string `json:"source"`
	WebhookURL string `json:"webhookUrl"`
}

type GetAllWebhookSourcesResponse = []WebhookSourceResponse
