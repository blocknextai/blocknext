package auth

import (
	"context"

	"github.com/blocknextai/go-packages/rbac"
	"github.com/google/uuid"
)

type UserPermissionChecker interface {
	HasPermission(ctx context.Context, userID uuid.UUID, requiredPermission *rbac.Permission) (bool, error)
}
