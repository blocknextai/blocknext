package regeneratewebhooktoken

import (
	"github.com/google/uuid"
)

type RegenerateWebhookTokenResponse struct {
	ID           uuid.UUID `json:"id"`
	WebhookToken string    `json:"webhookToken"`
}
