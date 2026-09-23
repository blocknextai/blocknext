package contract

import (
	apikeysApplicationAPIKeys "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/application/apikeys"
	apikeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/domain/apikeys"
)

type Scope = apikeysDomainAPIKeys.Scope

const ScopeMCPInvoke = apikeysDomainAPIKeys.ScopeMCPInvoke
const ScopeWorkflowsTrigger = apikeysDomainAPIKeys.ScopeWorkflowsTrigger

type APIKeyService = apikeysApplicationAPIKeys.APIKeyService
