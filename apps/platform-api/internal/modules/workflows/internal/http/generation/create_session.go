package generation

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	generationUseCases "github.com/blocknextai/platform-api/internal/modules/workflows/internal/usecases/generation"
)

type CreateSessionRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	Title          string    `json:"title"`
}

func NewCreateSessionHandler(service *generationUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CreateSessionRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := service.CreateSession(c.RequestCtx(), &generationUseCases.CreateSessionCommand{
			OrganizationID: new(request.OrganizationID),
			UserID:         new(userID),
			Title:          request.Title,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result))
	}
}
