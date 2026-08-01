package verificationtokens

import (
	"strings"
	"time"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/blocknextai/platform-api/internal/account/domain"
	"github.com/google/uuid"
)

type VerificationToken struct {
	database.BaseEntity

	UserID    uuid.UUID
	Purpose   Purpose
	TokenHash string
	Email     string
	ExpiresAt time.Time
}

func NewVerificationToken(
	userID uuid.UUID,
	purpose Purpose,
	tokenHash string,
	email string,
	ttl time.Duration,
) (*VerificationToken, error) {
	utcNow := time.Now().UTC()
	normalizedEmail := domain.NormalizeEmail(email)

	token := &VerificationToken{
		BaseEntity: database.BaseEntity{
			ID:        bnuuid.NewV7(),
			CreatedAt: utcNow,
			UpdatedAt: utcNow,
			DeletedAt: nil,
		},
		UserID:    userID,
		Purpose:   purpose,
		TokenHash: tokenHash,
		Email:     normalizedEmail,
		ExpiresAt: utcNow.Add(ttl),
	}

	return token.validateThenReturn()
}

func (t *VerificationToken) Consume() (*VerificationToken, error) {
	if time.Now().UTC().After(t.ExpiresAt) {
		return nil, ErrVerificationTokenExpired
	}

	utcNow := time.Now().UTC()
	t.DeletedAt = &utcNow
	t.UpdatedAt = utcNow

	return t, nil
}

func (t *VerificationToken) validateThenReturn() (*VerificationToken, error) {
	if t.UserID == uuid.Nil {
		return nil, ErrUserIDIsRequired
	}

	if !t.Purpose.IsValid() {
		return nil, ErrInvalidPurpose
	}

	if strings.TrimSpace(t.TokenHash) == "" {
		return nil, ErrTokenHashIsRequired
	}

	if strings.TrimSpace(t.Email) == "" {
		return nil, ErrEmailIsRequired
	}

	return t, nil
}
