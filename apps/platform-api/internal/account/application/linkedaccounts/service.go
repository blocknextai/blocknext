package linkedaccounts

import (
	"context"

	"github.com/blocknextai/platform-api/internal/account/domain/linkedaccounts"
	"github.com/google/uuid"
)

type LinkedAccountService interface {
	GetByIdentifier(ctx context.Context, identifier string) (*linkedaccounts.LinkedAccount, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*linkedaccounts.LinkedAccount, error)
	GetAllByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]*linkedaccounts.LinkedAccount, error)
}

type linkedAccountService struct {
	linkedAccountRepository linkedaccounts.LinkedAccountRepository
}

func NewLinkedAccountService(linkedAccountRepository linkedaccounts.LinkedAccountRepository) LinkedAccountService {
	return &linkedAccountService{linkedAccountRepository: linkedAccountRepository}
}

func (s *linkedAccountService) GetByIdentifier(ctx context.Context, identifier string) (*linkedaccounts.LinkedAccount, error) {
	return s.linkedAccountRepository.GetByIdentifier(ctx, identifier)
}

func (s *linkedAccountService) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*linkedaccounts.LinkedAccount, error) {
	return s.linkedAccountRepository.GetAllByUserID(ctx, userID)
}

func (s *linkedAccountService) GetAllByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]*linkedaccounts.LinkedAccount, error) {
	return s.linkedAccountRepository.GetAllByUserIDs(ctx, userIDs)
}
