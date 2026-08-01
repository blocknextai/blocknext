package createusertoken

import (
	"context"

	"github.com/blocknextai/platform-api/internal/account/domain"
	"github.com/blocknextai/platform-api/internal/account/domain/usernonces"
)

type Request struct {
	AuthProvider domain.AuthProvider
	Payload      map[string]any
}

type Response struct {
	ProviderID  string
	Identifier  string
	DisplayName string
	Nonce       *usernonces.UserNonce // TODO: remove this
}

type AuthProviderRegistry interface {
	GetProvider(authProvider domain.AuthProvider) (AuthProvider, error)
	GetAllProviderKeys() []domain.AuthProvider
}

type AuthProvider interface {
	Validate(ctx context.Context, request Request) (*Response, error)
	GenerateOAuthURL(userNonce *usernonces.UserNonce) (string, error)
	BuildLoginMessage(nonce string) string
}
