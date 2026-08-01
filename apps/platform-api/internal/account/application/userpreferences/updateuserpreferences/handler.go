package updateuserpreferences

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

func (h *Handler) Handle(ctx context.Context, command *UpdateUserPreferencesCommand) (*UpdateUserPreferencesResponse, error) {
	existingPreference, err := h.userPreferenceRepository.GetByUserID(ctx, command.UserID)
	if err != nil {
		if !errors.Is(err, userpreferences.ErrPreferenceNotFound) {
			return nil, err
		}

		existingPreference, err = userpreferences.NewDefault(command.UserID)
		if err != nil {
			return nil, err
		}
	}

	themeMode := existingPreference.ThemeMode
	themeColor := existingPreference.ThemeColor
	language := existingPreference.Language

	if command.Theme != nil {
		if command.Theme.Mode != nil {
			themeMode = *command.Theme.Mode
		}
		if command.Theme.Color != nil {
			themeColor = *command.Theme.Color
		}
	}
	if command.Language != nil {
		language = *command.Language
	}

	updatedPreference, err := existingPreference.Update(themeMode, themeColor, language)
	if err != nil {
		return nil, err
	}

	if err := h.userPreferenceRepository.Upsert(ctx, updatedPreference); err != nil {
		return nil, err
	}

	return &UpdateUserPreferencesResponse{
		Theme: ThemeResponse{
			Color: updatedPreference.ThemeColor,
			Mode:  updatedPreference.ThemeMode,
		},
		Language: updatedPreference.Language,
	}, nil
}
