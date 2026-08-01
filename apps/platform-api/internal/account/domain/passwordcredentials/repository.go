package passwordcredentials

import (
	"context"

	"github.com/google/uuid"
)

type PasswordCredentialRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*PasswordCredential, error)
	Create(ctx context.Context, credential *PasswordCredential) error
	Update(ctx context.Context, credential *PasswordCredential) error
	Delete(ctx context.Context, credential *PasswordCredential) error
}
