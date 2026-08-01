package getscopes

import (
	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/apikeys/domain/apikeys"
)

func MapGetScopesQueryToGetScopesResponse(scopes []apiKeysDomainAPIKeys.Scope) *GetScopesResponse {
	response := make(GetScopesResponse, 0, len(scopes))
	for _, scope := range scopes {
		response = append(response, ScopeResponse{
			Key: scope,
		})
	}
	return &response
}
