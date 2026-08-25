package oauth2

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/credentialoauth/application/oauth2/authurl"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type OrganizationAuthRequest struct {
	OrganizationID uuid.UUID `uri:"organizationId"`
	CredentialID   uuid.UUID `json:"credentialId"`
}

func NewOrganizationAuthHandler(handler cqrs.Handler[*authurl.AuthURLCommand, *authurl.AuthURLResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(OrganizationAuthRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &authurl.AuthURLCommand{
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
