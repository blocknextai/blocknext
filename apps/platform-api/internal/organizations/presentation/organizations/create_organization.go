package organizations

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/organizations/application/organizations/createorganization"
	"github.com/gofiber/fiber/v3"
)

type CreateOrganizationRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

func NewCreateOrganizationHandler(handler cqrs.Handler[*createorganization.CreateOrganizationCommand, *createorganization.CreateOrganizationResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CreateOrganizationRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &createorganization.CreateOrganizationCommand{
			UserID:      userID,
			Title:       request.Title,
			Description: request.Description,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result, resultPkg.WithMessage("organization created")))
	}
}
