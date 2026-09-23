package apikeys

import (
	"context"
	"time"

	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/domain/apikeys"
)

type GetAllAPIKeysQuery struct {
	OwnerType  commonDomain.OwnerType
	OwnerID    uuid.UUID
	Search     resultPkg.SearchRequest
	Pagination resultPkg.PaginationRequest
}

type APIKeyResponse struct {
	ID         uuid.UUID                   `json:"id"`
	Name       string                      `json:"name"`
	Scopes     apiKeysDomainAPIKeys.Scopes `json:"scopes"`
	LastUsedAt *time.Time                  `json:"lastUsedAt,omitempty"`
	CreatedAt  time.Time                   `json:"createdAt"`
	UpdatedAt  time.Time                   `json:"updatedAt"`
}

type GetAllAPIKeysResponse struct {
	Items      []*APIKeyResponse
	TotalCount int64
}

func (s *Service) GetAllAPIKeys(ctx context.Context, request *GetAllAPIKeysQuery) (*GetAllAPIKeysResponse, error) {
	apiKeys, totalCount, err := s.apiKeyRepository.GetAllByOwner(
		ctx,
		request.OwnerType,
		request.OwnerID,
		request.Search.Query,
		request.Pagination.Offset,
		request.Pagination.Limit,
	)
	if err != nil {
		return nil, err
	}

	return &GetAllAPIKeysResponse{
		Items:      MapAPIKeysToResponse(apiKeys),
		TotalCount: totalCount,
	}, nil
}

func MapAPIKeysToResponse(apiKeys []*apiKeysDomainAPIKeys.APIKey) []*APIKeyResponse {
	responses := make([]*APIKeyResponse, 0, len(apiKeys))
	for _, k := range apiKeys {
		responses = append(responses, &APIKeyResponse{
			ID:         k.ID,
			Name:       k.Name,
			Scopes:     k.Scopes,
			LastUsedAt: k.LastUsedAt,
			CreatedAt:  k.CreatedAt,
			UpdatedAt:  k.UpdatedAt,
		})
	}
	return responses
}
