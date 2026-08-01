package updatetrigger

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/secretmanager"
	triggersApplicationTriggers "github.com/blocknextai/platform-api/internal/triggers/application/triggers"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/triggers/domain/triggers"
)

type Handler struct {
	triggerRepository  triggersDomainTriggers.TriggerRepository
	secretManager      secretmanager.SecretManager
	transactionManager database.TransactionManager
}

func New(
	triggerRepository triggersDomainTriggers.TriggerRepository,
	secretManager secretmanager.SecretManager,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		triggerRepository:  triggerRepository,
		secretManager:      secretManager,
		transactionManager: transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *UpdateTriggerCommand) (*UpdateTriggerResponse, error) {
	var response *UpdateTriggerResponse
	err := h.transactionManager.ExecuteInTransaction(ctx, func(ctx context.Context) error {
		trigger, err := h.triggerRepository.GetByIDAndOrganizationID(ctx, command.TriggerID, command.OrganizationID)
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
				webhookSecret, err = h.resolveWebhookSecret(*command.WebhookSecret)
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

		err = h.triggerRepository.Update(ctx, updatedTrigger)
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

func (h *Handler) resolveWebhookSecret(plaintext string) (*string, error) {
	if plaintext == "" {
		return nil, nil
	}

	encrypted, err := h.secretManager.Encrypt(plaintext)
	if err != nil {
		return nil, triggersApplicationTriggers.ErrFailedToEncryptSecret
	}
	return &encrypted, nil
}
