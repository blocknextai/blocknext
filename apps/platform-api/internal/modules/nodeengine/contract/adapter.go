package contract

import (
	nodeengineApplicationAdapters "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/adapters"
	nodeengineDomainAdapters "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/adapters"
)

var Adapt = nodeengineDomainAdapters.Adapt
var AsWebhookVerifier = nodeengineDomainAdapters.AsWebhookVerifier
var ErrAdapterNotFound = nodeengineDomainAdapters.ErrAdapterNotFound
var ErrInvalidSignature = nodeengineDomainAdapters.ErrInvalidSignature
var GetAdapter = nodeengineDomainAdapters.GetAdapter

type TriggerContext = nodeengineDomainAdapters.TriggerContext
type VerificationRequest = nodeengineDomainAdapters.VerificationRequest
type VerificationResponse = nodeengineDomainAdapters.VerificationResponse

type AdapterService = nodeengineApplicationAdapters.AdapterService
