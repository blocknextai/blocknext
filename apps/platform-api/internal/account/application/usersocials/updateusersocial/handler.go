package updateusersocial

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/url/platform"
	"github.com/blocknextai/platform-api/internal/account/domain/usersocials"
)

type Handler struct {
	userSocialRepository usersocials.UserSocialRepository
	transactionManager   database.TransactionManager
}

func New(
	userSocialRepository usersocials.UserSocialRepository,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		userSocialRepository: userSocialRepository,
		transactionManager:   transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *UpdateUserSocialCommand) (*UpdateUserSocialResponse, error) {
	var response *UpdateUserSocialResponse
	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		socials, err := h.userSocialRepository.GetAllByUserID(txCtx, command.UserID)
		if err != nil {
			return err
		}

		for _, social := range socials {
			social, err = social.Delete()
			if err != nil {
				return err
			}

			err = h.userSocialRepository.Delete(txCtx, social)
			if err != nil {
				return err
			}
		}

		for index, item := range command.Items {
			platformName := platform.Detect(item.URL)
			social, err := usersocials.NewUserSocial(
				command.UserID,
				platformName,
				item.URL,
				index,
			)
			if err != nil {
				return err
			}

			err = h.userSocialRepository.Create(txCtx, social)
			if err != nil {
				return err
			}
		}

		response = &UpdateUserSocialResponse{}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
