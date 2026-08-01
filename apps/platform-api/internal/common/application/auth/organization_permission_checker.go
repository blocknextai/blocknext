package auth

import (
	"context"

	"github.com/blocknextai/go-packages/rbac"
	"github.com/google/uuid"
)

type OrganizationPermissionChecker interface {
	HasPermission(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID, requiredPermission *rbac.Permission) (bool, error)
}
