package linkedaccounts

import (
	"strings"
	"time"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/blocknextai/platform-api/internal/account/domain"
	"github.com/google/uuid"
)

type LinkedAccount struct {
	database.BaseEntity

	UserID       uuid.UUID
	AuthProvider domain.AuthProvider
	ProviderID   string
	Identifier   string
	DisplayName  *string
	IsPrimary    bool
	IsVerified   bool
}

func NewLinkedAccount(
	userID uuid.UUID,
	authProvider domain.AuthProvider,
	providerId string,
	identifier string,
	displayName *string,
	isPrimary bool,
	isVerified bool,
) (*LinkedAccount, error) {
	utcNow := time.Now().UTC()

	linkedAccount := &LinkedAccount{
		BaseEntity: database.BaseEntity{
			ID:        bnuuid.NewV7(),
			CreatedAt: utcNow,
			UpdatedAt: utcNow,
			DeletedAt: nil,
		},
		UserID:       userID,
		AuthProvider: authProvider,
		ProviderID:   providerId,
		Identifier:   identifier,
		DisplayName:  displayName,
		IsPrimary:    isPrimary,
		IsVerified:   isVerified,
	}

	return linkedAccount.validateThenReturn()
}

func (la *LinkedAccount) Update(
	identifier string,
	displayName *string,
) (*LinkedAccount, error) {
	la.UpdatedAt = time.Now().UTC()

	la.Identifier = identifier
	la.DisplayName = displayName

	return la.validateThenReturn()
}

func (la *LinkedAccount) MarkVerified() (*LinkedAccount, error) {
	la.IsVerified = true
	la.UpdatedAt = time.Now().UTC()

	return la.validateThenReturn()
}

func (la *LinkedAccount) ChangeEmailAndVerify(newEmail string) (*LinkedAccount, error) {
	normalized := domain.NormalizeEmail(newEmail)

	la.ProviderID = normalized
	la.Identifier = normalized
	la.DisplayName = new(normalized)
	la.IsVerified = true
	la.UpdatedAt = time.Now().UTC()

	return la.validateThenReturn()
}

func (la *LinkedAccount) Delete() (*LinkedAccount, error) {
	utcNow := time.Now().UTC()

	la.UpdatedAt = utcNow
	la.DeletedAt = new(utcNow)
	return la.validateThenReturn()
}

func (la *LinkedAccount) validateThenReturn() (*LinkedAccount, error) {
	if la.UserID == uuid.Nil {
		return nil, ErrUserIDIsRequired
	}

	if strings.TrimSpace(la.ProviderID) == "" {
		return nil, ErrProviderIDIsRequired
	}

	if strings.TrimSpace(la.Identifier) == "" {
		return nil, ErrIdentifierIsRequired
	}

	return la, nil
}
