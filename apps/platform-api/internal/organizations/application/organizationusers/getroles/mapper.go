package getroles

import (
	"github.com/blocknextai/go-packages/rbac"
)

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
