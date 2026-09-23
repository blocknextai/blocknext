package authorizationrequests

import (
	"context"

	"github.com/google/uuid"
)

type AuthorizationRequestRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*AuthorizationRequest, error)
	Create(ctx context.Context, request *AuthorizationRequest) error
	Update(ctx context.Context, request *AuthorizationRequest) error
}
