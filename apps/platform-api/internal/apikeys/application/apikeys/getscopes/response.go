package getscopes

import (
	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/apikeys/domain/apikeys"
)

type ScopeResponse struct {
	Key apiKeysDomainAPIKeys.Scope `json:"key"`
}

type GetScopesResponse = []ScopeResponse
