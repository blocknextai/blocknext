package toolinvocations

import (
	"context"

	"github.com/google/uuid"
)

type ToolInvocationRepository interface {
	GetAllByOrganizationID(
		ctx context.Context,
		organizationID uuid.UUID,
		searchQuery string,
		offset int,
		limit int,
	) ([]*ToolInvocation, int64, error)
	GetByIDAndOrganizationID(ctx context.Context, id uuid.UUID, organizationID uuid.UUID) (*ToolInvocation, error)

	Create(ctx context.Context, toolInvocation *ToolInvocation) error
}
