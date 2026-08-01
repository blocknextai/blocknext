package webhooks

import (
	"context"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	nodeEngineDomainAdapters "github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/triggers/domain/triggers"
	"github.com/google/uuid"
)

type Request struct {
	WebhookToken        string
	Source              string
	Payload             map[string]any
	VerificationRequest *nodeEngineDomainAdapters.VerificationRequest
}

type ResolvedWebhook struct {
	Verification *nodeEngineDomainAdapters.VerificationResponse

	TriggeredByUserID *uuid.UUID
	OrganizationID    uuid.UUID
	ExecutionContext  commonDomain.ExecutionContext
	ContextItemID     uuid.UUID
	RuntimeConfig     *triggersDomainTriggers.RuntimeConfig
	TriggerContext    *nodeEngineDomainAdapters.TriggerContext
}

type WebhookResolver interface {
	Resolve(ctx context.Context, req *Request) (*ResolvedWebhook, error)
}
