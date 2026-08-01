// Package errorhandler provides a Fiber error handler that converts errors into
// structured JSON responses based on their apperror kind.
package errorhandler

import (
	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/result"
	"github.com/gofiber/fiber/v3"
)

var (
	errInternalServer = apperror.Internal("internal server error")
)

// New returns a Fiber error handler that maps errors to JSON failure responses
// with the appropriate HTTP status. When isProduction is false, the underlying
// cause is included in the response.
func New(isProduction bool) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		appErr, ok := apperror.As(err)
		if !ok {
			appErr = errInternalServer.WithCause(err)
		}

		response := result.Fail[any](appErr)
		if appErr.Code != "" {
			response = response.WithCode(appErr.Code)
		}
		if !isProduction && appErr.Cause != nil {
			response = response.WithCause(appErr.Cause.Error())
		}
		return c.Status(apperror.KindToStatus(appErr.Kind)).JSON(response)
	}
}
