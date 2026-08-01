package presentation

import (
	"log/slog"

	uploadApplicationUpload "github.com/blocknextai/file-gateway-api/internal/upload/application/upload"
	uploadDomain "github.com/blocknextai/file-gateway-api/internal/upload/domain"
	"github.com/blocknextai/go-packages/result"
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router, service uploadApplicationUpload.Service) {
	router.Get("/upload/:uploadId", handleGetUploadRule(service))
	router.Post("/upload/:uploadId", handleUpload(service))
}

type UploadRequest struct {
	UploadID string `uri:"uploadId"`
}

func handleGetUploadRule(service uploadApplicationUpload.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UploadRequest)
		if err := c.Bind().All(request); err != nil {
			return uploadDomain.ErrInvalidUploadID
		}
		rule, err := service.GetRule(request.UploadID)
		if err != nil {
			return err
		}
		return c.JSON(result.Ok[any](rule))
	}
}

func handleUpload(service uploadApplicationUpload.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(UploadRequest)
		if err := c.Bind().All(request); err != nil {
			return uploadDomain.ErrInvalidUploadID
		}

		rule, err := service.GetRule(request.UploadID)
		if err != nil {
			return err
		}

		file, err := parseUploadForm(c, rule)
		if err != nil {
			return err
		}
		defer func() {
			if cerr := file.Close(); cerr != nil {
				slog.Warn("file close failed", "component", "UploadRoutes", "error", cerr)
			}
		}()

		uploadResult, err := service.Upload(c.RequestCtx(), rule, file)
		if err != nil {
			return err
		}

		return c.JSON(result.Ok(uploadResult, result.WithMessage("file uploaded")))
	}
}
