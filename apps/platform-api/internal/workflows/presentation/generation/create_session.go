package generation

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/workflows/application/generation/sessions/createsession"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CreateSessionRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	Title          string    `json:"title"`
}

func NewCreateSessionHandler(handler cqrs.Handler[*createsession.CreateSessionCommand, *createsession.CreateSessionResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(CreateSessionRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &createsession.CreateSessionCommand{
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
