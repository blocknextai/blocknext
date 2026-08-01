package oauth2

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/credentialoauth/application/oauth2/exchangecode"
	"github.com/gofiber/fiber/v3"
)

type CallbackRequest struct {
	Code  string `query:"code"`
	State string `query:"state"`
}

func NewCallbackHandler(handler cqrs.Handler[*exchangecode.ExchangeCodeQuery, *exchangecode.ExchangeCodeResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CallbackRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &exchangecode.ExchangeCodeQuery{
			Code:  request.Code,
			State: request.State,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("oauth connected")))
	}
}
