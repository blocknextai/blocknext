package userpreferences

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain/userpreferences"
)

type GetUserPreferencesQuery struct {
	UserID uuid.UUID
}

type Theme struct {
	Color string `json:"color"`
	Mode  string `json:"mode"`
}

type GetUserPreferencesResponse struct {
	Theme    Theme  `json:"theme"`
	Language string `json:"language"`
}

func (s *Service) GetUserPreferences(ctx context.Context, request *GetUserPreferencesQuery) (*GetUserPreferencesResponse, error) {
	pref, err := s.userPreferenceRepository.GetByUserID(ctx, request.UserID)
	if err != nil {
		if !errors.Is(err, userpreferences.ErrPreferenceNotFound) {
			return nil, err
		}
	}

	return MapUserPreferenceToResponse(pref), nil
}

func MapUserPreferenceToResponse(preference *userpreferences.UserPreference) *GetUserPreferencesResponse {
	if preference == nil {
		return &GetUserPreferencesResponse{
			Theme: Theme{
				Color: userpreferences.DefaultThemeColor,
				Mode:  userpreferences.DefaultThemeMode,
			},
			Language: userpreferences.DefaultLanguage,
		}
	}

	return &GetUserPreferencesResponse{
		Theme: Theme{
			Color: preference.ThemeColor,
			Mode:  preference.ThemeMode,
		},
		Language: preference.Language,
	}
}
