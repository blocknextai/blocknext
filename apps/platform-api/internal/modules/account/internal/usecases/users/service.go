package users

import (
	accountDomainUsers "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/users"
)

type Service struct {
	userRepository accountDomainUsers.UserRepository
}

func NewService(
	userRepository accountDomainUsers.UserRepository,
) *Service {
	return &Service{
		userRepository: userRepository,
	}
}
