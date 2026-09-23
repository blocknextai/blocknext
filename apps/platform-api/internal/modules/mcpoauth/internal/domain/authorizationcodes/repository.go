package authorizationcodes

import (
	"context"
)

type AuthorizationCodeRepository interface {
	GetByCodeHash(ctx context.Context, codeHash string) (*AuthorizationCode, error)
	Create(ctx context.Context, code *AuthorizationCode) error
	Update(ctx context.Context, code *AuthorizationCode) error
}
