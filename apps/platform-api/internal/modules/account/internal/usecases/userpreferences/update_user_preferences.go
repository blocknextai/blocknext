package userpreferences

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain/userpreferences"
)

type ThemeCommand struct {
	Color *string
	Mode  *string
}

type UpdateUserPreferencesCommand struct {
	UserID   uuid.UUID
	Theme    *ThemeCommand
	Language *string
}

type ThemeResponse struct {
	Color string `json:"color"`
	Mode  string `json:"mode"`
}

type UpdateUserPreferencesResponse struct {
	Theme    ThemeResponse `json:"theme"`
	Language string        `json:"language"`
}

func (s *Service) UpdateUserPreferences(ctx context.Context, command *UpdateUserPreferencesCommand) (*UpdateUserPreferencesResponse, error) {
	existingPreference, err := s.userPreferenceRepository.GetByUserID(ctx, command.UserID)
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

	if err := s.userPreferenceRepository.Upsert(ctx, updatedPreference); err != nil {
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
