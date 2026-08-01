// Package recovery provides a Fiber middleware that recovers from panics,
// logs them with a stack trace, and converts them into an internal server error.
package recovery

import (
	"errors"
	"log/slog"
	"runtime/debug"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/cast"
	"github.com/gofiber/fiber/v3"
)

var (
	errInternalServer = apperror.Internal("internal server error")
)

// New returns a Fiber handler that recovers from panics in downstream handlers,
// logs the panic with a stack trace, and returns an internal server error.
func New() fiber.Handler {
	return func(c fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Recovered from panic",
					"component", "recovery_middleware",
					"panic", r,
					"path", c.Path(),
					"method", c.Method(),
					"stack", string(debug.Stack()))

				err = errInternalServer.WithCause(errors.New(cast.ToString(r)))
			}
		}()
		return c.Next()
	}
}
