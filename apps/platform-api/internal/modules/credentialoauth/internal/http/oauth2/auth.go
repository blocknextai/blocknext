package oauth2

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	resultPkg "github.com/blocknextai/go-packages/result"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/http"
	oauth2UseCases "github.com/blocknextai/platform-api/internal/modules/credentialoauth/internal/usecases/oauth2"
)

type OrganizationAuthRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	CredentialID   uuid.UUID `json:"credentialId"`
}

func NewOrganizationAuthHandler(service *oauth2UseCases.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(OrganizationAuthRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := service.AuthURL(c.RequestCtx(), &oauth2UseCases.AuthURLCommand{
			OwnerType:    commonDomain.OwnerTypeOrganization,
			OwnerID:      request.OrganizationID,
			CredentialID: request.CredentialID,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(resultPkg.Ok(result, resultPkg.WithMessage("oauth url created")))
	}
}
