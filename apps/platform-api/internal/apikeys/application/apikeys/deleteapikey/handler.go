package deleteapikey

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

func (h *Handler) Handle(ctx context.Context, command *DeleteAPIKeyCommand) (*DeleteAPIKeyResponse, error) {
	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		apiKey, err := h.apiKeyRepository.GetByOwnerAndID(txCtx, command.OwnerType, command.OwnerID, command.APIKeyID)
		if err != nil {
			return apiKeysDomainAPIKeys.ErrAPIKeyNotFound
		}

		deletedAPIKey, err := apiKey.Delete()
		if err != nil {
			return err
		}

		return h.apiKeyRepository.Delete(txCtx, deletedAPIKey)
	})

	if err != nil {
		return nil, err
	}

	return &DeleteAPIKeyResponse{
		ID: command.APIKeyID,
	}, nil
}
