package getroles

import (
	"context"

	"github.com/blocknextai/go-packages/rbac"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(ctx context.Context, request *GetRolesQuery) (*GetRolesResponse, error) {
	roles := rbac.UserRoles

	return MapGetRolesQueryToGetRolesResponse(roles), nil
}
