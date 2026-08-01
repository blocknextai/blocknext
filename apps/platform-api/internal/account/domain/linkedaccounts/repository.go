package linkedaccounts

import (
	"context"

	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	"github.com/google/uuid"
)

type LinkedAccountRepository interface {
	GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*LinkedAccount, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*LinkedAccount, error)
	GetByAuthProviderAndProviderID(ctx context.Context, authProvider accountDomain.AuthProvider, providerID string) (*LinkedAccount, error)
	GetByAuthProviderAndUserID(ctx context.Context, authProvider accountDomain.AuthProvider, userID uuid.UUID) (*LinkedAccount, error)
	GetByProviderID(ctx context.Context, providerID string) (*LinkedAccount, error)
	GetByIdentifier(ctx context.Context, identifier string) (*LinkedAccount, error)
	Create(ctx context.Context, linkedAccount *LinkedAccount) error
	Update(ctx context.Context, linkedAccount *LinkedAccount) error
	Delete(ctx context.Context, linkedAccount *LinkedAccount) error
	GetAllByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]*LinkedAccount, error)
}
