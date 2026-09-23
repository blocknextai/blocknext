package webhooks

import (
	"context"

	"github.com/google/uuid"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	triggersDomainTriggers "github.com/blocknextai/platform-api/internal/modules/triggers/internal/domain/triggers"
)

type Request struct {
	WebhookToken        string
	Source              string
	Payload             map[string]any
	VerificationRequest *nodeengineContract.VerificationRequest
}

type ResolvedWebhook struct {
	Verification *nodeengineContract.VerificationResponse

	TriggeredByUserID *uuid.UUID
	OrganizationID    uuid.UUID
	ExecutionContext  commonDomain.ExecutionContext
	ContextItemID     uuid.UUID
	RuntimeConfig     *triggersDomainTriggers.RuntimeConfig
	TriggerContext    *nodeengineContract.TriggerContext
}

type WebhookResolver interface {
	Resolve(ctx context.Context, req *Request) (*ResolvedWebhook, error)
}
