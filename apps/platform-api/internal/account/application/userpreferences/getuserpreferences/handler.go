package getuserpreferences

import (
	"context"
	"errors"

	"github.com/blocknextai/platform-api/internal/account/domain/userpreferences"
)

type Handler struct {
	userPreferenceRepository userpreferences.UserPreferenceRepository
}

func New(
	userPreferenceRepository userpreferences.UserPreferenceRepository,
) *Handler {
	return &Handler{
		userPreferenceRepository: userPreferenceRepository,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetUserPreferencesQuery) (*GetUserPreferencesResponse, error) {
	pref, err := h.userPreferenceRepository.GetByUserID(ctx, request.UserID)
	if err != nil {
		if !errors.Is(err, userpreferences.ErrPreferenceNotFound) {
			return nil, err
		}
	}

	return MapUserPreferenceToResponse(pref), nil
}
