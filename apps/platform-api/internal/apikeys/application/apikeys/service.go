package apikeys

import (
	"context"

	apiKeysDomain "github.com/blocknextai/platform-api/internal/apikeys/domain/apikeys"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/google/uuid"
)

type APIKeyService interface {
	GetByOwnerAndID(
		ctx context.Context,
		ownerType commonDomain.OwnerType,
		ownerID uuid.UUID,
		id uuid.UUID,
	) (*apiKeysDomain.APIKey, error)
}

type apiKeyService struct {
	repository apiKeysDomain.APIKeyRepository
}

func NewAPIKeyService(repository apiKeysDomain.APIKeyRepository) APIKeyService {
	return &apiKeyService{
		repository: repository,
	}
}

func (s *apiKeyService) GetByOwnerAndID(
	ctx context.Context,
	ownerType commonDomain.OwnerType,
	ownerID uuid.UUID,
	id uuid.UUID,
) (*apiKeysDomain.APIKey, error) {
	return s.repository.GetByOwnerAndID(ctx, ownerType, ownerID, id)
}
