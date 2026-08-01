package createapikey

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

func (h *Handler) Handle(ctx context.Context, command *CreateAPIKeyCommand) (*CreateAPIKeyResponse, error) {
	generated, err := apiKeysDomainAPIKeys.GenerateKey()
	if err != nil {
		return nil, err
	}

	var apiKey *apiKeysDomainAPIKeys.APIKey

	err = h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		newKey, err := apiKeysDomainAPIKeys.New(
			command.OwnerType,
			command.OwnerID,
			command.Name,
			generated.Hash,
			command.Scopes,
		)
		if err != nil {
			return err
		}

		if err := h.apiKeyRepository.Create(txCtx, newKey); err != nil {
			return err
		}

		apiKey = newKey
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &CreateAPIKeyResponse{
		ID:   apiKey.ID,
		Name: apiKey.Name,
		Key:  generated.Plain,
	}, nil
}
