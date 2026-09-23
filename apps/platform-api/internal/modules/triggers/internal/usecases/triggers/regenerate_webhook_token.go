package triggers

import (
	"context"

	"github.com/google/uuid"

	triggersApplicationTriggers "github.com/blocknextai/platform-api/internal/modules/triggers/internal/application/triggers"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/modules/triggers/internal/domain/triggers"
)

type RegenerateWebhookTokenCommand struct {
	OrganizationID uuid.UUID
	TriggerID      uuid.UUID
}

type RegenerateWebhookTokenResponse struct {
	ID           uuid.UUID `json:"id"`
	WebhookToken string    `json:"webhookToken"`
}

func (s *Service) RegenerateWebhookToken(ctx context.Context, command *RegenerateWebhookTokenCommand) (*RegenerateWebhookTokenResponse, error) {
	var triggerID uuid.UUID
	var webhookToken string

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		existing, err := s.triggerRepository.GetByIDAndOrganizationID(txCtx, command.TriggerID, command.OrganizationID)
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

		if err := s.triggerRepository.Update(txCtx, regenerated); err != nil {
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
