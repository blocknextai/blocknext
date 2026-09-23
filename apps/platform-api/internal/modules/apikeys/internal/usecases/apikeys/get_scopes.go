package apikeys

import (
	"context"

	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/domain/apikeys"
)

type GetScopesQuery struct{}

type ScopeResponse struct {
	Key apiKeysDomainAPIKeys.Scope `json:"key"`
}

type GetScopesResponse = []ScopeResponse

func (s *Service) GetScopes(ctx context.Context, query *GetScopesQuery) (*GetScopesResponse, error) {
	scopes := make([]apiKeysDomainAPIKeys.Scope, 0, len(apiKeysDomainAPIKeys.AllScopes))
	for scope := range apiKeysDomainAPIKeys.AllScopes {
		scopes = append(scopes, scope)
	}

	return MapGetScopesQueryToGetScopesResponse(scopes), nil
}

func MapGetScopesQueryToGetScopesResponse(scopes []apiKeysDomainAPIKeys.Scope) *GetScopesResponse {
	response := make(GetScopesResponse, 0, len(scopes))
	for _, scope := range scopes {
		response = append(response, ScopeResponse{
			Key: scope,
		})
	}
	return &response
}
