package usersocials

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/usersocials/updateusersocial"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type UpdateUserSocialItem struct {
	URL string `json:"url"`
}

type UpdateUserSocialRequest struct {
	Items []UpdateUserSocialItem `json:"items"`
}

func NewUpdateUserSocialHandler(handler cqrs.Handler[*updateusersocial.UpdateUserSocialCommand, *updateusersocial.UpdateUserSocialResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UpdateUserSocialRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		items := make([]updateusersocial.UpdateUserSocialItem, 0, len(request.Items))
		for _, item := range request.Items {
			items = append(items, updateusersocial.UpdateUserSocialItem{
				URL: item.URL,
			})
		}

		result, err := handler.Handle(c.RequestCtx(), &updateusersocial.UpdateUserSocialCommand{
			UserID: userID,
			Items:  items,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("socials updated")))
	}
}
