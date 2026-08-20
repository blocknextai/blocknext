package triggers

import (
	"time"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type Trigger struct {
	database.BaseEntity

	OrganizationID    uuid.UUID
	TriggeredByUserID *uuid.UUID
	ExecutionContext  domain.ExecutionContext
	ContextItemID     uuid.UUID
	Type              TriggerType
	CronPattern       *string
	Timezone          *string
	WebhookTokenHash  *string
	WebhookSecret     *string
	RuntimeConfig     *RuntimeConfig
	IsActive          bool
}

func New(
	organizationID uuid.UUID,
	triggeredByUserID *uuid.UUID,
	executionContext domain.ExecutionContext,
	contextItemID uuid.UUID,
	triggerType TriggerType,
	cronPattern *string,
	timezone *string,
	webhookTokenHash *string,
	webhookSecret *string,
	runtimeConfig *RuntimeConfig,
	isActive bool,
) (*Trigger, error) {
	utcNow := time.Now().UTC()

	flowTrigger := &Trigger{
		ID:                bnuuid.NewV7(),
		CreatedAt:         utcNow,
		UpdatedAt:         utcNow,
		OrganizationID:    organizationID,
		TriggeredByUserID: triggeredByUserID,
		ExecutionContext:  executionContext,
		ContextItemID:     contextItemID,
		Type:              triggerType,
		CronPattern:       cronPattern,
		Timezone:          timezone,
		WebhookTokenHash:  webhookTokenHash,
		WebhookSecret:     webhookSecret,
		RuntimeConfig:     runtimeConfig,
		IsActive:          isActive,
	}

	return flowTrigger.validateThenReturn()
}

func (ft *Trigger) update(
	isActive bool,
	runtimeConfig *RuntimeConfig,
) (*Trigger, error) {
	ft.UpdatedAt = time.Now().UTC()

	ft.IsActive = isActive
	ft.RuntimeConfig = runtimeConfig

	return ft.validateThenReturn()
}

func (ft *Trigger) UpdateSchedule(
	isActive bool,
	cronPattern *string,
	timezone *string,
	runtimeConfig *RuntimeConfig,
) (*Trigger, error) {
	if ft.Type != TriggerTypeSchedule {
		return nil, ErrNotScheduleTrigger
	}

	ft.CronPattern = cronPattern
	ft.Timezone = timezone

	return ft.update(isActive, runtimeConfig)
}

func (ft *Trigger) UpdateWebhook(
	isActive bool,
	webhookTokenHash *string,
	webhookSecret *string,
	runtimeConfig *RuntimeConfig,
) (*Trigger, error) {
	if ft.Type != TriggerTypeWebhook {
		return nil, ErrNotWebhookTrigger
	}

	ft.WebhookTokenHash = webhookTokenHash
	ft.WebhookSecret = webhookSecret

	return ft.update(isActive, runtimeConfig)
}

func (ft *Trigger) Delete() (*Trigger, error) {
	utcNow := time.Now().UTC()

	ft.UpdatedAt = utcNow
	ft.DeletedAt = new(utcNow)

	return ft.validateThenReturn()
}

func (ft *Trigger) validateThenReturn() (*Trigger, error) {
	return ft, nil
}
