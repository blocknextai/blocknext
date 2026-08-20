package userpreferences

import (
	"time"

	"github.com/blocknextai/go-packages/database"
	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/google/uuid"
)

const (
	DefaultThemeMode  = "system"
	DefaultThemeColor = "default"
	DefaultLanguage   = "en"
)

type UserPreference struct {
	database.BaseEntity

	UserID     uuid.UUID
	ThemeMode  string
	ThemeColor string
	Language   string
}

func NewDefault(userID uuid.UUID) (*UserPreference, error) {
	return NewUserPreference(userID, DefaultThemeMode, DefaultThemeColor, DefaultLanguage)
}

func NewUserPreference(
	userID uuid.UUID,
	themeMode string,
	themeColor string,
	language string,
) (*UserPreference, error) {
	utcNow := time.Now().UTC()

	pref := &UserPreference{
		ID:         bnuuid.NewV7(),
		CreatedAt:  utcNow,
		UpdatedAt:  utcNow,
		DeletedAt:  nil,
		UserID:     userID,
		ThemeMode:  themeMode,
		ThemeColor: themeColor,
		Language:   language,
	}

	return pref.validateThenReturn()
}

func (p *UserPreference) Update(
	themeMode string,
	themeColor string,
	language string,
) (*UserPreference, error) {
	p.UpdatedAt = time.Now().UTC()

	p.ThemeMode = themeMode
	p.ThemeColor = themeColor
	p.Language = language

	return p.validateThenReturn()
}

func (p *UserPreference) validateThenReturn() (*UserPreference, error) {
	if p.UserID == uuid.Nil {
		return nil, ErrPreferenceNotFound
	}

	return p, nil
}
