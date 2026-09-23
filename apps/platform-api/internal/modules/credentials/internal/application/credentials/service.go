package credentials

import (
	"context"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/secretmanager"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	"github.com/blocknextai/platform-api/internal/modules/credentials/internal/domain/credentials"
)

type CredentialService interface {
	GetAllByOwner(ctx context.Context, ownerType commonDomain.OwnerType, ownerID uuid.UUID, searchQuery string, offset int, limit int) ([]*credentials.Credential, int64, error)
	GetByIDForOwner(ctx context.Context, credentialID uuid.UUID, ownerType commonDomain.OwnerType, ownerID uuid.UUID) (*credentials.CredentialInfo, error)
	SaveCredentialForOwner(ctx context.Context, credentialID uuid.UUID, ownerType commonDomain.OwnerType, ownerID uuid.UUID, data any) error
}

type credentialService struct {
	credentialRepository credentials.CredentialRepository
	secretManager        secretmanager.SecretManager
}

func NewCredentialService(
	credentialRepository credentials.CredentialRepository,
	secretManager secretmanager.SecretManager,
) CredentialService {
	return &credentialService{
		credentialRepository: credentialRepository,
		secretManager:        secretManager,
	}
}

func (s *credentialService) GetAllByOwner(
	ctx context.Context,
	ownerType commonDomain.OwnerType,
	ownerID uuid.UUID,
	searchQuery string,
	offset int,
	limit int,
) ([]*credentials.Credential, int64, error) {
	return s.credentialRepository.GetAllByOwner(ctx, ownerType, ownerID, searchQuery, offset, limit)
}

func (s *credentialService) GetByIDForOwner(ctx context.Context, credentialID uuid.UUID, ownerType commonDomain.OwnerType, ownerID uuid.UUID) (*credentials.CredentialInfo, error) {
	credential, err := s.credentialRepository.GetByIDAndOwner(ctx, credentialID, ownerType, ownerID)
	if err != nil {
		return nil, err
	}

	var decryptedCredential map[string]any
	if err := s.secretManager.Decrypt(credential.Data, &decryptedCredential); err != nil {
		return nil, err
	}

	return &credentials.CredentialInfo{
		Key:        credential.Key,
		SourceType: credential.SourceType,
		Data:       decryptedCredential,
	}, nil
}

func (s *credentialService) SaveCredentialForOwner(ctx context.Context, credentialID uuid.UUID, ownerType commonDomain.OwnerType, ownerID uuid.UUID, data any) error {
	credential, err := s.credentialRepository.GetByIDAndOwner(ctx, credentialID, ownerType, ownerID)
	if err != nil {
		return err
	}

	encryptedData, err := s.secretManager.Encrypt(data)
	if err != nil {
		return err
	}

	credential.Data = encryptedData

	return s.credentialRepository.Update(ctx, credential)
}
