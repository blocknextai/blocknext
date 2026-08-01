package passwordcredentials

import (
	"strings"
	"time"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/google/uuid"
)

type PasswordCredential struct {
	database.BaseEntity

	UserID       uuid.UUID
	PasswordHash string
}

func NewPasswordCredential(
	userID uuid.UUID,
	passwordHash string,
) (*PasswordCredential, error) {
	utcNow := time.Now().UTC()

	credential := &PasswordCredential{
		BaseEntity: database.BaseEntity{
			ID:        bnuuid.NewV7(),
			CreatedAt: utcNow,
			UpdatedAt: utcNow,
			DeletedAt: nil,
		},
		UserID:       userID,
		PasswordHash: passwordHash,
	}

	return credential.validateThenReturn()
}

func (pc *PasswordCredential) Delete() (*PasswordCredential, error) {
	utcNow := time.Now().UTC()

	pc.UpdatedAt = utcNow
	pc.DeletedAt = &utcNow

	return pc.validateThenReturn()
}

func (pc *PasswordCredential) ChangePassword(passwordHash string) (*PasswordCredential, error) {
	pc.PasswordHash = passwordHash
	pc.UpdatedAt = time.Now().UTC()

	return pc.validateThenReturn()
}

func (pc *PasswordCredential) validateThenReturn() (*PasswordCredential, error) {
	if pc.UserID == uuid.Nil {
		return nil, ErrPasswordCredentialNotFound
	}

	if strings.TrimSpace(pc.PasswordHash) == "" {
		return nil, ErrPasswordHashIsRequired
	}

	return pc, nil
}
