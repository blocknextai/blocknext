package auth

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/cache"
	"github.com/blocknextai/go-packages/rbac"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	"github.com/blocknextai/platform-api/internal/common/application/auth"
	"github.com/google/uuid"
)

const (
	userPermissionCacheTTL    = 5 * time.Minute
	userPermissionCachePrefix = "auth:permission:user:"
)

type UserPermissionChecker struct {
	userRepository accountDomainUsers.UserRepository
	cacheService   cache.Service
}

func NewUserPermissionChecker(
	userRepository accountDomainUsers.UserRepository,
	cacheService cache.Service,
) auth.UserPermissionChecker {
	return &UserPermissionChecker{
		userRepository: userRepository,
		cacheService:   cacheService,
	}
}

func (c *UserPermissionChecker) HasPermission(ctx context.Context, userID uuid.UUID, requiredPermission *rbac.Permission) (bool, error) {
	role, err := c.getRole(ctx, userID)
	if err != nil {
		if errors.Is(err, accountDomainUsers.ErrUserNotFound) {
			return false, nil
		}
		return false, err
	}

	if role == "" {
		return false, nil
	}

	return rbac.HasUserPermission(role, requiredPermission.Code), nil
}

func (c *UserPermissionChecker) getRole(ctx context.Context, userID uuid.UUID) (string, error) {
	cacheKey := c.buildCacheKey(userID)

	cachedRole, err := c.cacheService.Get(ctx, cacheKey)
	if err != nil {
		slog.Error("failed to get user permission from cache",
			"userId", userID,
			"error", err,
		)
	}

	if cachedRole != "" {
		return cachedRole, nil
	}

	user, err := c.userRepository.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	if err := c.cacheService.Set(ctx, cacheKey, user.Role, userPermissionCacheTTL); err != nil {
		slog.Error("failed to set user permission to cache",
			"userId", userID,
			"error", err,
		)
	}

	return user.Role, nil
}

func (c *UserPermissionChecker) buildCacheKey(userID uuid.UUID) string {
	userIDStr := userID.String()
	var builder strings.Builder
	builder.Grow(len(userPermissionCachePrefix) + len(userIDStr))
	builder.WriteString(userPermissionCachePrefix)
	builder.WriteString(userIDStr)
	return builder.String()
}
