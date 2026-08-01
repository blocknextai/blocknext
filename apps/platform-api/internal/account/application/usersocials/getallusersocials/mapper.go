package getallusersocials

import (
	"github.com/blocknextai/platform-api/internal/account/domain/usersocials"
)

func MapToResponseList(socials []*usersocials.UserSocial) *GetAllUserSocialsResponse {
	response := make(GetAllUserSocialsResponse, 0, len(socials))
	for _, social := range socials {
		response = append(response, UserSocialResponse{
			Platform:  social.Platform,
			URL:       social.URL,
			SortOrder: social.SortOrder,
		})
	}
	return &response
}
