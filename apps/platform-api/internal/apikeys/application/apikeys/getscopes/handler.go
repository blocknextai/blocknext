package getscopes

import (
	"context"

	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/apikeys/domain/apikeys"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, query *GetScopesQuery) (*GetScopesResponse, error) {
	scopes := make([]apiKeysDomainAPIKeys.Scope, 0, len(apiKeysDomainAPIKeys.AllScopes))
	for scope := range apiKeysDomainAPIKeys.AllScopes {
		scopes = append(scopes, scope)
	}

	return MapGetScopesQueryToGetScopesResponse(scopes), nil
}
