package apikeys

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/apikeys/application/apikeys/deleteapikey"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
)

type DeleteOrganizationAPIKeyRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	APIKeyID       uuid.UUID `uri:"apiKeyId"`
}

func NewDeleteOrganizationAPIKeyHandler(handler cqrs.Handler[*deleteapikey.DeleteAPIKeyCommand, *deleteapikey.DeleteAPIKeyResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DeleteOrganizationAPIKeyRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &deleteapikey.DeleteAPIKeyCommand{
			OwnerType: commonDomain.OwnerTypeOrganization,
			OwnerID:   request.OrganizationID,
			APIKeyID:  request.APIKeyID,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("api key deleted")))
	}
}
