package usersocials

import (
	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain/usersocials"
)

type Service struct {
	userSocialRepository usersocials.UserSocialRepository
	transactionManager   database.TransactionManager
}

func NewService(
	userSocialRepository usersocials.UserSocialRepository,
	transactionManager database.TransactionManager,
) *Service {
	return &Service{
		userSocialRepository: userSocialRepository,
		transactionManager:   transactionManager,
	}
}
