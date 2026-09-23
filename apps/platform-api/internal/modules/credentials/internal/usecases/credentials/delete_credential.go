package credentials

import (
	"context"

	"github.com/google/uuid"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/modules/credentials/internal/application/credentials"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/modules/credentials/internal/domain/credentials"
)

type DeleteCredentialCommand struct {
	OwnerType    commonDomain.OwnerType
	OwnerID      uuid.UUID
	CredentialID uuid.UUID
}

type DeleteCredentialResponse struct{}

func (s *Service) DeleteCredential(ctx context.Context, request *DeleteCredentialCommand) (*DeleteCredentialResponse, error) {
	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		credential, err := s.credentialRepository.GetByIDAndOwner(txCtx, request.CredentialID, request.OwnerType, request.OwnerID)
		if err != nil {
			return credentialsDomainCredentials.ErrCredentialNotFound
		}

		credential, err = credential.Delete()
		if err != nil {
			return err
		}

		err = s.credentialRepository.Delete(txCtx, credential)
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
