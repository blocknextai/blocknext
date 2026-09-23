package oauth2

import (
	"github.com/gofiber/fiber/v3"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	oauth2UseCases "github.com/blocknextai/platform-api/internal/modules/credentialoauth/internal/usecases/oauth2"
)

type CallbackRequest struct {
	Code  string `query:"code"`
	State string `query:"state"`
}

func NewCallbackHandler(service *oauth2UseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CallbackRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.ExchangeCode(c.RequestCtx(), &oauth2UseCases.ExchangeCodeQuery{
			Code:  request.Code,
			State: request.State,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("oauth connected")))
	}
}
