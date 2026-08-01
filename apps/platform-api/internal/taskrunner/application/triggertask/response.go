package triggertask

import (
	"github.com/google/uuid"
)

type TriggerTaskResponse struct {
	ID           uuid.UUID `json:"id"`
	WebhookToken *string   `json:"webhookToken,omitempty"`
}
