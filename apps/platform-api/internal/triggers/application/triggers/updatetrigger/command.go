package updatetrigger

import (
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/triggers/domain/triggers"
	"github.com/google/uuid"
)

type UpdateTriggerCommand struct {
	OrganizationID uuid.UUID
	TriggerID      uuid.UUID
	IsActive       *bool
	CronPattern    *string
	Timezone       *string
	WebhookSecret  *string
	RuntimeConfig  *triggersDomainTriggers.RuntimeConfig
}
