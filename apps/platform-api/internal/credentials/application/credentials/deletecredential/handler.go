package deletecredential

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/credentials/application/credentials"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/credentials/domain/credentials"
)

type Handler struct {
	credentialRepository credentialsDomainCredentials.CredentialRepository
	transactionManager   database.TransactionManager
}

func New(
	credentialRepository credentialsDomainCredentials.CredentialRepository,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		credentialRepository: credentialRepository,
		transactionManager:   transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, request *DeleteCredentialCommand) (*DeleteCredentialResponse, error) {
	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		credential, err := h.credentialRepository.GetByIDAndOwner(txCtx, request.CredentialID, request.OwnerType, request.OwnerID)
		if err != nil {
			return credentialsDomainCredentials.ErrCredentialNotFound
		}

		credential, err = credential.Delete()
		if err != nil {
			return err
		}

		err = h.credentialRepository.Delete(txCtx, credential)
		if err != nil {
			return credentialsApplicationCredentials.ErrFailedToDeleteCredential.WithCause(err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &DeleteCredentialResponse{}, nil
}
