package apikeys

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/apikeys/application/apikeys/regenerateapikey"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
)

type RegenerateOrganizationAPIKeyRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	APIKeyID       uuid.UUID `uri:"apiKeyId"`
}

func NewRegenerateOrganizationAPIKeyHandler(handler cqrs.Handler[*regenerateapikey.RegenerateAPIKeyCommand, *regenerateapikey.RegenerateAPIKeyResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(RegenerateOrganizationAPIKeyRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &regenerateapikey.RegenerateAPIKeyCommand{
			OwnerType: commonDomain.OwnerTypeOrganization,
			OwnerID:   request.OrganizationID,
			APIKeyID:  request.APIKeyID,
		})
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("api key regenerated")))
	}
}
