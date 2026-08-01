package webhooks

import (
	"context"
	"log/slog"

	"github.com/blocknextai/go-packages/secretmanager"
	nodeEngineDomainAdapters "github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
	triggersApplicationTriggers "github.com/blocknextai/platform-api/internal/triggers/application/triggers"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/triggers/domain/triggers"
	triggersDomainWebhooks "github.com/blocknextai/platform-api/internal/triggers/domain/webhooks"
)

type webhookResolver struct {
	triggerRepository triggersDomainTriggers.TriggerRepository
	secretManager     secretmanager.SecretManager
}

func NewWebhookResolver(
	triggerRepository triggersDomainTriggers.TriggerRepository,
	secretManager secretmanager.SecretManager,
) WebhookResolver {
	return &webhookResolver{
		triggerRepository: triggerRepository,
		secretManager:     secretManager,
	}
}

func (h *webhookResolver) Resolve(ctx context.Context, req *Request) (*ResolvedWebhook, error) {
	tokenHash := triggersApplicationTriggers.HashWebhookToken(req.WebhookToken)

	trigger, err := h.triggerRepository.GetByWebhookTokenHash(ctx, tokenHash)
	if err != nil {
		slog.ErrorContext(ctx, "Webhook trigger not found",
			"component", "webhook.trigger",
			"source", req.Source,
			"error", err,
		)
		return nil, triggersDomainWebhooks.ErrTriggerNotFound
	}

	if !trigger.IsActive {
		return nil, triggersDomainWebhooks.ErrTriggerInactive
	}

	adapter, ok := nodeEngineDomainAdapters.GetAdapter(req.Source)
	if !ok {
		return nil, nodeEngineDomainAdapters.ErrAdapterNotFound
	}

	if verifier, ok := nodeEngineDomainAdapters.AsWebhookVerifier(adapter); ok && req.VerificationRequest != nil {
		secret, err := h.resolveSecret(trigger.WebhookSecret)
		if err != nil {
			return nil, err
		}

		response, err := verifier.Verify(req.VerificationRequest, secret)
		if err != nil {
			return nil, err
		}
		if response != nil {
			return &ResolvedWebhook{Verification: response}, nil
		}

		if err := verifier.ValidateSignature(req.VerificationRequest, secret); err != nil {
			return nil, nodeEngineDomainAdapters.ErrInvalidSignature
		}
	}

	triggerContext, err := nodeEngineDomainAdapters.Adapt(req.Source, req.Payload)
	if err != nil {
		return nil, err
	}

	return &ResolvedWebhook{
		TriggeredByUserID: trigger.TriggeredByUserID,
		OrganizationID:    trigger.OrganizationID,
		ExecutionContext:  trigger.ExecutionContext,
		ContextItemID:     trigger.ContextItemID,
		RuntimeConfig:     trigger.RuntimeConfig,
		TriggerContext:    triggerContext,
	}, nil
}

func (h *webhookResolver) resolveSecret(encrypted *string) (*string, error) {
	if encrypted == nil {
		return nil, nil
	}

	var plaintext string
	if err := h.secretManager.Decrypt(*encrypted, &plaintext); err != nil {
		return nil, triggersDomainWebhooks.ErrFailedToDecryptSecret
	}
	return &plaintext, nil
}
