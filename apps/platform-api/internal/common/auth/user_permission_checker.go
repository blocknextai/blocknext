package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/rbac"
)

type UserPermissionChecker interface {
	HasPermission(ctx context.Context, userID uuid.UUID, requiredPermission *rbac.Permission) (bool, error)
}
