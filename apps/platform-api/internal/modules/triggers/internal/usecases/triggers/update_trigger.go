package triggers

import (
	"context"

	"github.com/google/uuid"

	triggersApplicationTriggers "github.com/blocknextai/platform-api/internal/modules/triggers/internal/application/triggers"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/modules/triggers/internal/domain/triggers"
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

type UpdateTriggerResponse struct {
	ID uuid.UUID `json:"id"`
}

func (s *Service) UpdateTrigger(ctx context.Context, command *UpdateTriggerCommand) (*UpdateTriggerResponse, error) {
	var response *UpdateTriggerResponse
	err := s.transactionManager.ExecuteInTransaction(ctx, func(ctx context.Context) error {
		trigger, err := s.triggerRepository.GetByIDAndOrganizationID(ctx, command.TriggerID, command.OrganizationID)
		if err != nil {
			return err
		}

		isActive := trigger.IsActive
		if command.IsActive != nil {
			isActive = *command.IsActive
		}

		runtimeConfig := trigger.RuntimeConfig
		if command.RuntimeConfig != nil {
			runtimeConfig = command.RuntimeConfig
		}

		var updatedTrigger *triggersDomainTriggers.Trigger
		switch trigger.Type {
		case triggersDomainTriggers.TriggerTypeSchedule:
			cronPattern := trigger.CronPattern
			if command.CronPattern != nil {
				cronPattern = command.CronPattern
			}

			timezone := trigger.Timezone
			if command.Timezone != nil {
				timezone = command.Timezone
			}

			updatedTrigger, err = trigger.UpdateSchedule(isActive, cronPattern, timezone, runtimeConfig)
		case triggersDomainTriggers.TriggerTypeWebhook:
			webhookSecret := trigger.WebhookSecret
			if command.WebhookSecret != nil {
				webhookSecret, err = s.resolveWebhookSecret(*command.WebhookSecret)
				if err != nil {
					return err
				}
			}

			updatedTrigger, err = trigger.UpdateWebhook(isActive, trigger.WebhookTokenHash, webhookSecret, runtimeConfig)
		default:
			return triggersDomainTriggers.ErrUnsupportedType
		}
		if err != nil {
			return err
		}

		err = s.triggerRepository.Update(ctx, updatedTrigger)
		if err != nil {
			return err
		}

		response = &UpdateTriggerResponse{
			ID: updatedTrigger.ID,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) resolveWebhookSecret(plaintext string) (*string, error) {
	if plaintext == "" {
		return nil, nil
	}

	encrypted, err := s.secretManager.Encrypt(plaintext)
	if err != nil {
		return nil, triggersApplicationTriggers.ErrFailedToEncryptSecret
	}
	return &encrypted, nil
}
