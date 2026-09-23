package apikeys

import (
	"context"

	"github.com/google/uuid"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/domain/apikeys"
)

type RegenerateAPIKeyCommand struct {
	OwnerType commonDomain.OwnerType
	OwnerID   uuid.UUID
	APIKeyID  uuid.UUID
}

type RegenerateAPIKeyResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Key  string    `json:"key"`
}

func (s *Service) RegenerateAPIKey(ctx context.Context, command *RegenerateAPIKeyCommand) (*RegenerateAPIKeyResponse, error) {
	generated, err := apiKeysDomainAPIKeys.GenerateKey()
	if err != nil {
		return nil, err
	}

	var apiKey *apiKeysDomainAPIKeys.APIKey

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		existing, err := s.apiKeyRepository.GetByOwnerAndID(txCtx, command.OwnerType, command.OwnerID, command.APIKeyID)
		if err != nil {
			return apiKeysDomainAPIKeys.ErrAPIKeyNotFound
		}

		regenerated, err := existing.Regenerate(generated.Hash)
		if err != nil {
			return err
		}

		if err := s.apiKeyRepository.Update(txCtx, regenerated); err != nil {
			return err
		}

		apiKey = regenerated

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &RegenerateAPIKeyResponse{
		ID:   apiKey.ID,
		Name: apiKey.Name,
		Key:  generated.Plain,
	}, nil
}
