package userpreferences

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/userpreferences/updateuserpreferences"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type ThemeRequest struct {
	Color *string `json:"color"`
	Mode  *string `json:"mode"`
}

type UpdateUserPreferencesRequest struct {
	Theme    *ThemeRequest `json:"theme"`
	Language *string       `json:"language"`
}

func NewUpdateUserPreferencesHandler(handler cqrs.Handler[*updateuserpreferences.UpdateUserPreferencesCommand, *updateuserpreferences.UpdateUserPreferencesResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UpdateUserPreferencesRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		command := &updateuserpreferences.UpdateUserPreferencesCommand{
			UserID:   userID,
			Language: request.Language,
		}
		if request.Theme != nil {
			command.Theme = &updateuserpreferences.ThemeCommand{
				Color: request.Theme.Color,
				Mode:  request.Theme.Mode,
			}
		}

		result, err := handler.Handle(c.RequestCtx(), command)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
