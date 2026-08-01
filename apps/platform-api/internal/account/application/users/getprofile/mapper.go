package getprofile

import (
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
)

func MapUserToResponse(user *accountDomainUsers.User, permissions []string) *GetProfileResponse {
	return &GetProfileResponse{
		Role:        user.Role,
		Permissions: permissions,
		IsVerified:  user.IsVerified,
	}
}
