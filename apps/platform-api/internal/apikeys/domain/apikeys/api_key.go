package apikeys

import (
	"strings"
	"time"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type APIKey struct {
	database.BaseEntity

	OwnerType  commonDomain.OwnerType
	OwnerID    uuid.UUID
	Name       string
	KeyHash    string
	Scopes     Scopes
	LastUsedAt *time.Time
}

func New(
	ownerType commonDomain.OwnerType,
	ownerID uuid.UUID,
	name string,
	keyHash string,
	scopes Scopes,
) (*APIKey, error) {
	now := time.Now().UTC()

	apiKey := &APIKey{
		ID:        bnuuid.NewV7(),
		CreatedAt: now,
		UpdatedAt: now,
		OwnerType: ownerType,
		OwnerID:   ownerID,
		Name:      name,
		KeyHash:   keyHash,
		Scopes:    scopes,
	}

	return apiKey.validateThenReturn()
}

func (k *APIKey) HasScope(scope Scope) bool {
	return k.Scopes.Has(scope)
}

func (k *APIKey) Update(
	lastUsedAt *time.Time,
) (*APIKey, error) {
	k.LastUsedAt = lastUsedAt
	k.UpdatedAt = time.Now().UTC()

	return k.validateThenReturn()
}

func (k *APIKey) Regenerate(keyHash string) (*APIKey, error) {
	k.KeyHash = keyHash

	return k.Update(nil)
}

func (k *APIKey) Delete() (*APIKey, error) {
	now := time.Now().UTC()

	k.UpdatedAt = now
	k.DeletedAt = new(now)

	return k.validateThenReturn()
}

func (k *APIKey) validateThenReturn() (*APIKey, error) {
	if strings.TrimSpace(k.Name) == "" {
		return nil, ErrInvalidAPIKeyName
	}

	if !k.OwnerType.IsValid() {
		return nil, ErrInvalidOwnerType
	}

	if k.OwnerID == uuid.Nil {
		return nil, ErrInvalidOwnerID
	}

	if k.KeyHash == "" {
		return nil, ErrInvalidAPIKey
	}

	if len(k.Scopes) == 0 {
		return nil, ErrInvalidAPIKeyScopes
	}
	for _, scope := range k.Scopes {
		if !scope.IsValid() {
			return nil, ErrInvalidAPIKeyScopes
		}
	}

	return k, nil
}
