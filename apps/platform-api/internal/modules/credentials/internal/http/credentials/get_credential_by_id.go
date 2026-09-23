package credentials

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	credentialsUseCases "github.com/blocknextai/platform-api/internal/modules/credentials/internal/usecases/credentials"
)

type GetOrganizationCredentialByIDRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	CredentialID   uuid.UUID `uri:"credentialId"`
}

func NewGetOrganizationCredentialByIDHandler(service *credentialsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(GetOrganizationCredentialByIDRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.GetCredentialByID(c.RequestCtx(), &credentialsUseCases.GetCredentialByIDQuery{
			OwnerType:    commonDomain.OwnerTypeOrganization,
			OwnerID:      request.OrganizationID,
			CredentialID: request.CredentialID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result))
	}
}
