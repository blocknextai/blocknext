package credentials

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonHTTP "github.com/blocknextai/platform-api/internal/common/presentation/http"
	"github.com/blocknextai/platform-api/internal/credentials/application/credentials/updatecredential"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UpdateUserCredentialRequest struct {
	CredentialID uuid.UUID      `uri:"credentialId"`
	Key          string         `json:"key"`
	Title        string         `json:"title"`
	Data         map[string]any `json:"data"`
}

func NewUpdateUserCredentialHandler(handler cqrs.Handler[*updatecredential.UpdateCredentialCommand, *updatecredential.UpdateCredentialResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UpdateUserCredentialRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		userID := commonHTTP.GetUserID(c)

		result, err := handler.Handle(c.RequestCtx(), &updatecredential.UpdateCredentialCommand{
			ID:        request.CredentialID,
			OwnerType: commonDomain.OwnerTypeUser,
			OwnerID:   userID,
			Key:       request.Key,
			Title:     request.Title,
			Data:      request.Data,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("credential updated")))
	}
}

type UpdateOrganizationCredentialRequest struct {
	OrganizationID uuid.UUID      `uri:"organizationId"`
	CredentialID   uuid.UUID      `uri:"credentialId"`
	Key            string         `json:"key"`
	Title          string         `json:"title"`
	Data           map[string]any `json:"data"`
}

func NewUpdateOrganizationCredentialHandler(handler cqrs.Handler[*updatecredential.UpdateCredentialCommand, *updatecredential.UpdateCredentialResponse]) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UpdateOrganizationCredentialRequest)
		if err := c.Bind().All(request); err != nil {
			return commonHTTP.ErrInvalidRequest
		}

		result, err := handler.Handle(c.RequestCtx(), &updatecredential.UpdateCredentialCommand{
			ID:        request.CredentialID,
			OwnerType: commonDomain.OwnerTypeOrganization,
			OwnerID:   request.OrganizationID,
			Key:       request.Key,
			Title:     request.Title,
			Data:      request.Data,
		})

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(result, resultPkg.WithMessage("credential updated")))
	}
}
