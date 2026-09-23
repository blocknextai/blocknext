package usersocials

import (
	"context"

	"github.com/google/uuid"

	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain/usersocials"
)

type GetAllUserSocialsQuery struct {
	UserID uuid.UUID
}

type UserSocialResponse struct {
	Platform  string `json:"platform"`
	URL       string `json:"url"`
	SortOrder int    `json:"sortOrder"`
}

type GetAllUserSocialsResponse = []UserSocialResponse

func (s *Service) GetAllUserSocials(ctx context.Context, request *GetAllUserSocialsQuery) (*GetAllUserSocialsResponse, error) {
	socials, err := s.userSocialRepository.GetAllByUserID(ctx, request.UserID)
	if err != nil {
		return nil, err
	}

	return MapToResponseList(socials), nil
}

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
