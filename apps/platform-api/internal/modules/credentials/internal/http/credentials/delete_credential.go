package credentials

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	credentialsUseCases "github.com/blocknextai/platform-api/internal/modules/credentials/internal/usecases/credentials"
)

type DeleteOrganizationCredentialRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	CredentialID   uuid.UUID `uri:"credentialId"`
}

func NewDeleteOrganizationCredentialHandler(service *credentialsUseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DeleteOrganizationCredentialRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.DeleteCredential(c.RequestCtx(), &credentialsUseCases.DeleteCredentialCommand{
			OwnerType:    commonDomain.OwnerTypeOrganization,
			OwnerID:      request.OrganizationID,
			CredentialID: request.CredentialID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("credential deleted")))
	}
}
