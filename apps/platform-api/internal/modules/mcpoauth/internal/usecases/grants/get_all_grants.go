package grants

import (
	"context"
	"time"

	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	mcpOAuthDomainGrants "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/grants"
)

type GetAllGrantsQuery struct {
	UserID     uuid.UUID
	Search     resultPkg.SearchRequest
	Pagination resultPkg.PaginationRequest
}

type GrantResponse struct {
	ID             uuid.UUID `json:"id"`
	ClientID       string    `json:"clientId"`
	ClientName     string    `json:"clientName"`
	OrganizationID uuid.UUID `json:"organizationId"`
	Scopes         []string  `json:"scopes"`
	Resource       string    `json:"resource"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type GetAllGrantsResponse struct {
	Items      []*GrantResponse
	TotalCount int64
}

func (s *Service) GetAllGrants(ctx context.Context, query *GetAllGrantsQuery) (*GetAllGrantsResponse, error) {
	grants, totalCount, err := s.grantRepository.GetAllActiveByUserID(
		ctx,
		query.UserID,
		query.Search.Query,
		query.Pagination.Offset,
		query.Pagination.Limit,
	)
	if err != nil {
		return nil, err
	}

	clientNames := make(map[string]string, len(grants))
	for _, g := range grants {
		if _, ok := clientNames[g.ClientID]; ok {
			continue
		}

		client, err := s.clientResolver.Resolve(ctx, g.ClientID)
		if err != nil {
			continue
		}

		clientNames[g.ClientID] = client.Name
	}

	return &GetAllGrantsResponse{
		Items:      MapGrantsToResponse(grants, clientNames),
		TotalCount: totalCount,
	}, nil
}

func MapGrantsToResponse(grants []*mcpOAuthDomainGrants.Grant, clientNames map[string]string) []*GrantResponse {
	responses := make([]*GrantResponse, 0, len(grants))
	for _, g := range grants {
		clientName, ok := clientNames[g.ClientID]
		if !ok {
			clientName = g.ClientID
		}

		responses = append(responses, &GrantResponse{
			ID:             g.ID,
			ClientID:       g.ClientID,
			ClientName:     clientName,
			OrganizationID: g.OrganizationID,
			Scopes:         g.Scopes.Strings(),
			Resource:       g.Resource,
			CreatedAt:      g.CreatedAt,
			UpdatedAt:      g.UpdatedAt,
		})
	}
	return responses
}
