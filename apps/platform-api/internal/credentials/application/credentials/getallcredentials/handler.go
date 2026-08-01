package getallcredentials

import (
	"context"

	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/credentials/application/credentials"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/credentials/domain/credentials"
)

type Handler struct {
	credentialRepository credentialsDomainCredentials.CredentialRepository
}

func New(
	credentialRepository credentialsDomainCredentials.CredentialRepository,
) *Handler {
	return &Handler{
		credentialRepository: credentialRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetAllCredentialsQuery) (*GetAllCredentialsResponse, error) {
	credentials, total, err := h.credentialRepository.GetAllByOwner(
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

	return &GetAllCredentialsResponse{
		Items:      credentialsApplicationCredentials.MapCredentialsToResponse(credentials),
		TotalCount: total,
	}, nil
}
