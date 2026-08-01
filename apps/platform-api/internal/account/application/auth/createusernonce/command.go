package createusernonce

import (
	"github.com/blocknextai/platform-api/internal/account/domain"
)

type CreateUserNonceCommand struct {
	AuthProvider domain.AuthProvider
	ProviderID   *string
}
