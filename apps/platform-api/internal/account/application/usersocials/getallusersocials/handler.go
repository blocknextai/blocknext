package getallusersocials

import (
	"context"

	"github.com/blocknextai/platform-api/internal/account/domain/usersocials"
)

type Handler struct {
	userSocialRepository usersocials.UserSocialRepository
}

func New(
	userSocialRepository usersocials.UserSocialRepository,
) *Handler {
	return &Handler{
		userSocialRepository: userSocialRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetAllUserSocialsQuery) (*GetAllUserSocialsResponse, error) {
	socials, err := h.userSocialRepository.GetAllByUserID(ctx, request.UserID)
	if err != nil {
		return nil, err
	}

	return MapToResponseList(socials), nil
}
