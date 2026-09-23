package apikeys

import (
	"context"

	"github.com/google/uuid"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/domain/apikeys"
)

type DeleteAPIKeyCommand struct {
	OwnerType commonDomain.OwnerType
	OwnerID   uuid.UUID
	APIKeyID  uuid.UUID
}

type DeleteAPIKeyResponse struct {
	ID uuid.UUID `json:"id"`
}

func (s *Service) DeleteAPIKey(ctx context.Context, command *DeleteAPIKeyCommand) (*DeleteAPIKeyResponse, error) {
	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		apiKey, err := s.apiKeyRepository.GetByOwnerAndID(txCtx, command.OwnerType, command.OwnerID, command.APIKeyID)
		if err != nil {
			return apiKeysDomainAPIKeys.ErrAPIKeyNotFound
		}

		deletedAPIKey, err := apiKey.Delete()
		if err != nil {
			return err
		}

		return s.apiKeyRepository.Delete(txCtx, deletedAPIKey)
	})

	if err != nil {
		return nil, err
	}

	return &DeleteAPIKeyResponse{
		ID: command.APIKeyID,
	}, nil
}
