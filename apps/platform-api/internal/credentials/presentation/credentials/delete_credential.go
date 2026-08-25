package credentials

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/credentials/application/credentials/deletecredential"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type DeleteOrganizationCredentialRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	CredentialID   uuid.UUID `uri:"credentialId"`
}

func NewDeleteOrganizationCredentialHandler(handler cqrs.Handler[*deletecredential.DeleteCredentialCommand, *deletecredential.DeleteCredentialResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DeleteOrganizationCredentialRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &deletecredential.DeleteCredentialCommand{
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
