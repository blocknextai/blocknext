package linkedaccounts

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/account/application/linkedaccounts/getalllinkedaccounts"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/gofiber/fiber/v3"
)

func NewGetAllLinkedAccountsHandler(handler cqrs.Handler[*getalllinkedaccounts.GetAllLinkedAccountsQuery, *getalllinkedaccounts.GetAllLinkedAccountsResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &getalllinkedaccounts.GetAllLinkedAccountsQuery{
			UserID: userID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
