package userpreferences

import (
	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain/userpreferences"
)

type Service struct {
	userPreferenceRepository userpreferences.UserPreferenceRepository
}

func NewService(
	userPreferenceRepository userpreferences.UserPreferenceRepository,
) *Service {
	return &Service{
		userPreferenceRepository: userPreferenceRepository,
	}
}
