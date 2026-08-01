package regeneratewebhooktoken

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	triggersApplicationTriggers "github.com/blocknextai/platform-api/internal/triggers/application/triggers"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/triggers/domain/triggers"
	"github.com/google/uuid"
)

type Handler struct {
	triggerRepository  triggersDomainTriggers.TriggerRepository
	transactionManager database.TransactionManager
}

func New(
	triggerRepository triggersDomainTriggers.TriggerRepository,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		triggerRepository:  triggerRepository,
		transactionManager: transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *RegenerateWebhookTokenCommand) (*RegenerateWebhookTokenResponse, error) {
	var triggerID uuid.UUID
	var webhookToken string

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		existing, err := h.triggerRepository.GetByIDAndOrganizationID(txCtx, command.TriggerID, command.OrganizationID)
		if err != nil {
			return triggersDomainTriggers.ErrTriggerNotFound
		}

		plainToken, tokenHash, err := triggersApplicationTriggers.GenerateWebhookToken()
		if err != nil {
			return err
		}

		regenerated, err := existing.UpdateWebhook(existing.IsActive, new(tokenHash), existing.WebhookSecret, existing.RuntimeConfig)
		if err != nil {
			return err
		}

		if err := h.triggerRepository.Update(txCtx, regenerated); err != nil {
			return err
		}

		triggerID = regenerated.ID
		webhookToken = plainToken

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &RegenerateWebhookTokenResponse{
		ID:           triggerID,
		WebhookToken: webhookToken,
	}, nil
}
