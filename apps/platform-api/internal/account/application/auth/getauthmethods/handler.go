package getauthmethods

import (
	"context"

	"github.com/blocknextai/platform-api/internal/account/application/auth/createusertoken"
)

type Handler struct {
	authProviderRegistry createusertoken.AuthProviderRegistry
	passwordEnabled      bool
	magicLinkEnabled     bool
}

func New(authProviderRegistry createusertoken.AuthProviderRegistry, passwordEnabled bool, magicLinkEnabled bool) *Handler {
	return &Handler{
		authProviderRegistry: authProviderRegistry,
		passwordEnabled:      passwordEnabled,
		magicLinkEnabled:     magicLinkEnabled,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetAuthMethodsQuery) (*GetAuthMethodsResponse, error) {
	providers := h.authProviderRegistry.GetAllProviderKeys()
	return MapAuthMethodsToResponse(providers, h.passwordEnabled, h.magicLinkEnabled), nil
}
