package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/rbac"
)

type OrganizationPermissionChecker interface {
	HasPermission(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID, requiredPermission *rbac.Permission) (bool, error)
}
