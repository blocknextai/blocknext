package getuserpreferences

import (
	"github.com/blocknextai/platform-api/internal/account/domain/userpreferences"
)

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
