package auth

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/cache"
	"github.com/blocknextai/go-packages/rbac"
	"github.com/blocknextai/platform-api/internal/common/application/auth"
	"github.com/blocknextai/platform-api/internal/organizations/domain/organizationusers"
	"github.com/google/uuid"
)

const (
	organizationPermissionCacheTTL    = 5 * time.Minute
	organizationPermissionCachePrefix = "auth:permission:organization:"
)

type OrganizationPermissionChecker struct {
	organizationUserRepository organizationusers.OrganizationUserRepository
	cacheService               cache.Service
}

func NewOrganizationPermissionChecker(
	organizationUserRepository organizationusers.OrganizationUserRepository,
	cacheService cache.Service,
) auth.OrganizationPermissionChecker {
	return &OrganizationPermissionChecker{
		organizationUserRepository: organizationUserRepository,
		cacheService:               cacheService,
	}
}

func (c *OrganizationPermissionChecker) HasPermission(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID, requiredPermission *rbac.Permission) (bool, error) {
	role, err := c.getRole(ctx, organizationID, userID)
	if err != nil {
		if errors.Is(err, organizationusers.ErrOrganizationUserNotFound) {
			return false, nil
		}
		return false, err
	}

	if role == "" {
		return false, nil
	}

	return rbac.HasOrganizationPermission(role, requiredPermission.Code), nil
}

func (c *OrganizationPermissionChecker) getRole(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID) (string, error) {
	cacheKey := c.buildCacheKey(organizationID, userID)

	cachedRole, err := c.cacheService.Get(ctx, cacheKey)
	if err != nil {
		slog.Error("failed to get organization permission from cache",
			"organizationId", organizationID,
			"userId", userID,
			"error", err,
		)
	}

	if cachedRole != "" {
		return cachedRole, nil
	}

	organizationUser, err := c.organizationUserRepository.GetByOrganizationIDAndUserID(ctx, organizationID, userID)
	if err != nil {
		return "", err
	}

	if err := c.cacheService.Set(ctx, cacheKey, organizationUser.Role, organizationPermissionCacheTTL); err != nil {
		slog.Error("failed to set organization permission to cache",
			"organizationId", organizationID,
			"userId", userID,
			"error", err,
		)
	}

	return organizationUser.Role, nil
}

func (c *OrganizationPermissionChecker) buildCacheKey(organizationID uuid.UUID, userID uuid.UUID) string {
	organizationIDStr := organizationID.String()
	userIDStr := userID.String()
	var builder strings.Builder
	builder.Grow(len(organizationPermissionCachePrefix) + len(organizationIDStr) + 1 + len(userIDStr))
	builder.WriteString(organizationPermissionCachePrefix)
	builder.WriteString(organizationIDStr)
	builder.WriteString(":")
	builder.WriteString(userIDStr)
	return builder.String()
}
