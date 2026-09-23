package apikeys

import (
	"context"
	"strings"

	"github.com/google/uuid"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/domain/apikeys"
)

type CreateAPIKeyCommand struct {
	OwnerType commonDomain.OwnerType
	OwnerID   uuid.UUID
	Name      string
	Scopes    apiKeysDomainAPIKeys.Scopes
}

func (c *CreateAPIKeyCommand) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return apiKeysDomainAPIKeys.ErrInvalidAPIKeyName
	}

	if len(c.Scopes) == 0 {
		return apiKeysDomainAPIKeys.ErrInvalidAPIKeyScopes
	}
	for _, scope := range c.Scopes {
		if !scope.IsValid() {
			return apiKeysDomainAPIKeys.ErrInvalidAPIKeyScopes
		}
	}

	return nil
}

type CreateAPIKeyResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Key  string    `json:"key"`
}

func (s *Service) CreateAPIKey(ctx context.Context, command *CreateAPIKeyCommand) (*CreateAPIKeyResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	generated, err := apiKeysDomainAPIKeys.GenerateKey()
	if err != nil {
		return nil, err
	}

	var apiKey *apiKeysDomainAPIKeys.APIKey

	err = s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
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

		if err := s.apiKeyRepository.Create(txCtx, newKey); err != nil {
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
