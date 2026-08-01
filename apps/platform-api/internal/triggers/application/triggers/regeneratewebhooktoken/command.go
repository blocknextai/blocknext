package regeneratewebhooktoken

import (
	"github.com/google/uuid"
)

type RegenerateWebhookTokenCommand struct {
	OrganizationID uuid.UUID
	TriggerID      uuid.UUID
}
