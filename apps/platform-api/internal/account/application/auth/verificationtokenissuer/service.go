package verificationtokenissuer

import (
	"context"
	"time"

	accountApplicationAuth "github.com/blocknextai/platform-api/internal/account/application/auth"
	"github.com/blocknextai/platform-api/internal/account/application/auth/verificationtoken"
	accountDomainVerificationTokens "github.com/blocknextai/platform-api/internal/account/domain/verificationtokens"
	"github.com/google/uuid"
)

type Service interface {
	Issue(ctx context.Context, userID uuid.UUID, purpose accountDomainVerificationTokens.Purpose, email string, ttl time.Duration) (string, error)
}

type service struct {
	repository accountDomainVerificationTokens.VerificationTokenRepository
}

func NewService(repository accountDomainVerificationTokens.VerificationTokenRepository) Service {
	return &service{repository: repository}
}

func (s *service) Issue(ctx context.Context, userID uuid.UUID, purpose accountDomainVerificationTokens.Purpose, email string, ttl time.Duration) (string, error) {
	plainToken, tokenHash, err := verificationtoken.Generate()
	if err != nil {
		return "", accountApplicationAuth.ErrFailedToGenerateVerificationToken
	}

	if err := s.repository.InvalidateAllForUserAndPurpose(ctx, userID, purpose); err != nil {
		return "", err
	}

	token, err := accountDomainVerificationTokens.NewVerificationToken(userID, purpose, tokenHash, email, ttl)
	if err != nil {
		return "", err
	}

	if err := s.repository.Create(ctx, token); err != nil {
		return "", err
	}

	return plainToken, nil
}
