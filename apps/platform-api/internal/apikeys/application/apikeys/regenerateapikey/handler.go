package regenerateapikey

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/apikeys/domain/apikeys"
)

type Handler struct {
	apiKeyRepository   apiKeysDomainAPIKeys.APIKeyRepository
	transactionManager database.TransactionManager
}

func New(
	apiKeyRepository apiKeysDomainAPIKeys.APIKeyRepository,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		apiKeyRepository:   apiKeyRepository,
		transactionManager: transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *RegenerateAPIKeyCommand) (*RegenerateAPIKeyResponse, error) {
	generated, err := apiKeysDomainAPIKeys.GenerateKey()
	if err != nil {
		return nil, err
	}

	var apiKey *apiKeysDomainAPIKeys.APIKey

	err = h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		existing, err := h.apiKeyRepository.GetByOwnerAndID(txCtx, command.OwnerType, command.OwnerID, command.APIKeyID)
		if err != nil {
			return apiKeysDomainAPIKeys.ErrAPIKeyNotFound
		}

		regenerated, err := existing.Regenerate(generated.Hash)
		if err != nil {
			return err
		}

		if err := h.apiKeyRepository.Update(txCtx, regenerated); err != nil {
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
