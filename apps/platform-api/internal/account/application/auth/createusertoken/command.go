package createusertoken

import (
	"github.com/blocknextai/platform-api/internal/account/domain"
)

type CreateUserTokenCommand struct {
	AuthProvider domain.AuthProvider
	Payload      map[string]any
	IPAddress    string
	UserAgent    string
}
