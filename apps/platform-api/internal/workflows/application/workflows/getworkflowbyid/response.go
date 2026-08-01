package getworkflowbyid

import (
	"github.com/blocknextai/go-packages/dag"
	"github.com/google/uuid"
)

type LinkedAccount struct {
	AuthProvider string  `json:"authProvider"`
	DisplayName  *string `json:"displayName"`
	IsPrimary    bool    `json:"isPrimary"`
}

type Owner struct {
	ID             uuid.UUID       `json:"id"`
	Alias          string          `json:"alias"`
	IsVerified     bool            `json:"isVerified"`
	LinkedAccounts []LinkedAccount `json:"linkedAccounts"`
}

type GetWorkflowByIDResponse struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organizationId"`
	Owner          *Owner     `json:"owner"`
	Title          string     `json:"title"`
	Description    *string    `json:"description"`
	IsPinned       bool       `json:"isPinned"`
	Nodes          []dag.Node `json:"nodes"`
	Edges          []dag.Edge `json:"edges"`
}
