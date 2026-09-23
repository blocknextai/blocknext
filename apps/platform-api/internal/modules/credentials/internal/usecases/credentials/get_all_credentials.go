package credentials

import (
	"context"

	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/modules/credentials/internal/application/credentials"
)

type GetAllCredentialsQuery struct {
	OwnerType  commonDomain.OwnerType
	OwnerID    uuid.UUID
	Search     resultPkg.SearchRequest
	Pagination resultPkg.PaginationRequest
}

type GetAllCredentialsResponse struct {
	Items      []credentialsApplicationCredentials.CredentialResponse `json:"items"`
	TotalCount int64                                                  `json:"totalCount"`
}

func (s *Service) GetAllCredentials(ctx context.Context, request *GetAllCredentialsQuery) (*GetAllCredentialsResponse, error) {
	credentials, total, err := s.credentialRepository.GetAllByOwner(
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
