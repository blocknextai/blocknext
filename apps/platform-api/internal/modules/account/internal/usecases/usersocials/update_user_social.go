package usersocials

import (
	"context"
	"net/url"
	"strings"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/url/platform"
	accountApplicationUserSocials "github.com/blocknextai/platform-api/internal/modules/account/internal/application/usersocials"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain/usersocials"
)

type UpdateUserSocialItem struct {
	URL string
}

type UpdateUserSocialCommand struct {
	UserID uuid.UUID
	Items  []UpdateUserSocialItem
}

func (c *UpdateUserSocialCommand) Validate() error {
	for _, item := range c.Items {
		if strings.TrimSpace(item.URL) == "" {
			return accountApplicationUserSocials.ErrURLRequired
		}

		parsedURL, err := url.Parse(item.URL)
		if err != nil {
			return accountApplicationUserSocials.ErrInvalidURL
		}

		if parsedURL.Scheme != "https" {
			return accountApplicationUserSocials.ErrHttpsRequired
		}

		if parsedURL.Host == "" {
			return accountApplicationUserSocials.ErrInvalidURL
		}
	}

	return nil
}

type UpdateUserSocialResponse struct{}

func (s *Service) UpdateUserSocial(ctx context.Context, command *UpdateUserSocialCommand) (*UpdateUserSocialResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, err
	}

	var response *UpdateUserSocialResponse
	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		socials, err := s.userSocialRepository.GetAllByUserID(txCtx, command.UserID)
		if err != nil {
			return err
		}

		for _, social := range socials {
			social, err = social.Delete()
			if err != nil {
				return err
			}

			err = s.userSocialRepository.Delete(txCtx, social)
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

			err = s.userSocialRepository.Create(txCtx, social)
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
