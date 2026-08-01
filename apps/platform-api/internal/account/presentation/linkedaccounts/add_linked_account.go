package linkedaccounts

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/linkedaccounts/addlinkedaccount"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

type AddLinkedAccountRequest struct {
	AuthProvider accountDomain.AuthProvider `json:"authProvider"`
	Payload      map[string]any             `json:"payload"`
}

func NewAddLinkedAccountHandler(handler cqrs.Handler[*addlinkedaccount.AddLinkedAccountCommand, *addlinkedaccount.AddLinkedAccountResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(AddLinkedAccountRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &addlinkedaccount.AddLinkedAccountCommand{
			UserID:       userID,
			AuthProvider: request.AuthProvider,
			Payload:      request.Payload,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result, resultPkg.WithMessage("linked account added")))
	}
}
