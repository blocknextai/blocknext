package http

import (
	resultPkg "github.com/blocknextai/go-packages/result"
	"github.com/gofiber/fiber/v3"
)

func RespondPaginated[T any](c fiber.Ctx, items []T, totalCount int64, p resultPkg.PaginationRequest) error {
	pagination := resultPkg.NewPagination(totalCount, p.Offset, p.Limit)
	return c.Status(fiber.StatusOK).JSON(resultPkg.Ok(items, resultPkg.WithPagination(pagination)))
}
