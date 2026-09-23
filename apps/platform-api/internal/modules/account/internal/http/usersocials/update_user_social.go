package usersocials

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	usersocialsUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/usersocials"
)

type UpdateUserSocialItem struct {
	URL string `json:"url"`
}

type UpdateUserSocialRequest struct {
	Items []UpdateUserSocialItem `json:"items"`
}

func NewUpdateUserSocialHandler(service *usersocialsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UpdateUserSocialRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		items := make([]usersocialsUseCases.UpdateUserSocialItem, 0, len(request.Items))
		for _, item := range request.Items {
			items = append(items, usersocialsUseCases.UpdateUserSocialItem{
				URL: item.URL,
			})
		}

		result, err := service.UpdateUserSocial(c.RequestCtx(), &usersocialsUseCases.UpdateUserSocialCommand{
			UserID: userID,
			Items:  items,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("socials updated")))
	}
}
