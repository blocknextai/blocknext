package organizationusers

import (
	"context"

	"github.com/blocknextai/go-packages/rbac"
)

type GetRolesQuery struct{}

type RoleResponse struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type GetRolesResponse = []RoleResponse

func (s *Service) GetRoles(ctx context.Context, request *GetRolesQuery) (*GetRolesResponse, error) {
	roles := rbac.OrganizationRoles

	return MapGetRolesQueryToGetRolesResponse(roles), nil
}

func MapGetRolesQueryToGetRolesResponse(roles []*rbac.Role) *GetRolesResponse {
	response := make(GetRolesResponse, 0, len(roles))
	for _, role := range roles {
		response = append(response, RoleResponse{
			Name:  role.Name,
			Score: role.Score,
		})
	}
	return &response
}
