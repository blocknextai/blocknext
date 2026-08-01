package presentation

import (
	"strings"

	commonDomain "github.com/blocknextai/file-gateway-api/internal/common/domain"
	downloadDomain "github.com/blocknextai/file-gateway-api/internal/download/domain"
	downloadDomainDownload "github.com/blocknextai/file-gateway-api/internal/download/domain/download"
	"github.com/blocknextai/go-packages/cast"
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router, service downloadDomainDownload.Service) {
	router.Post("/download", handleDownload(service))
}

type DownloadRequest struct {
	URL string `json:"url"`
}

func handleDownload(service downloadDomainDownload.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(DownloadRequest)
		if err := c.Bind().All(request); err != nil {
			return commonDomain.ErrInvalidRequest
		}

		if strings.TrimSpace(request.URL) == "" {
			return downloadDomain.ErrInvalidURL
		}

		info, err := service.Download(c.RequestCtx(), request.URL)
		if err != nil {
			return err
		}

		setResponseHeaders(c, info)
		c.Response().SetBodyStream(info.BodyReader, int(info.Size))
		return nil
	}
}

func setResponseHeaders(c fiber.Ctx, info *downloadDomainDownload.FileInfo) {
	safeFilename := sanitizeHeaderValue(info.Filename)
	c.Set("Content-Type", info.ContentType)
	c.Set("Content-Disposition", `attachment; filename="`+safeFilename+`"`)
	if info.Size > 0 {
		c.Set("Content-Length", cast.ToString(info.Size))
	}
	c.Set("X-Filename", safeFilename)
	c.Set("X-Type", info.Type)
	c.Set("X-Ext", info.Ext)
}

func sanitizeHeaderValue(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}
