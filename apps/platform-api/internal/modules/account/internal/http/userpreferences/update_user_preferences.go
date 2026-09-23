package userpreferences

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	userpreferencesUseCases "github.com/blocknextai/platform-api/internal/modules/account/internal/usecases/userpreferences"
)

type ThemeRequest struct {
	Color *string `json:"color"`
	Mode  *string `json:"mode"`
}

type UpdateUserPreferencesRequest struct {
	Theme    *ThemeRequest `json:"theme"`
	Language *string       `json:"language"`
}

func NewUpdateUserPreferencesHandler(service *userpreferencesUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UpdateUserPreferencesRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		command := &userpreferencesUseCases.UpdateUserPreferencesCommand{
			UserID:   userID,
			Language: request.Language,
		}
		if request.Theme != nil {
			command.Theme = &userpreferencesUseCases.ThemeCommand{
				Color: request.Theme.Color,
				Mode:  request.Theme.Mode,
			}
		}

		result, err := service.UpdateUserPreferences(c.RequestCtx(), command)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
