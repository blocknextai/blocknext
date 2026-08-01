package getallapikeys

import (
	"context"

	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/apikeys/domain/apikeys"
)

type Handler struct {
	apiKeyRepository apiKeysDomainAPIKeys.APIKeyRepository
}

func New(
	apiKeyRepository apiKeysDomainAPIKeys.APIKeyRepository,
) *Handler {
	return &Handler{
		apiKeyRepository: apiKeyRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetAllAPIKeysQuery) (*GetAllAPIKeysResponse, error) {
	apiKeys, totalCount, err := h.apiKeyRepository.GetAllByOwner(
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
