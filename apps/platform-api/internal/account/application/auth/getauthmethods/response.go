package getauthmethods

import (
	"github.com/blocknextai/platform-api/internal/account/domain"
)

type GetAuthMethodsResponse struct {
	Crypto    []domain.AuthProvider `json:"crypto"`
	OAuth     []domain.AuthProvider `json:"oauth"`
	Password  bool                  `json:"password"`
	MagicLink bool                  `json:"magicLink"`
}
