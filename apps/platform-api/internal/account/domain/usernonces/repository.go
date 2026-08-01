package usernonces

import (
	"context"

	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
)

type UserNonceRepository interface {
	GetByNonceAndAuthProvider(ctx context.Context, nonce string, authProvider accountDomain.AuthProvider) (*UserNonce, error)
	Create(ctx context.Context, userNonce *UserNonce) error
	Delete(ctx context.Context, userNonce *UserNonce) error
}
