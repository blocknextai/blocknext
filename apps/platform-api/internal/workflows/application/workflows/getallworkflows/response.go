package getallworkflows

import (
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

type WorkflowResponse struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organizationId"`
	Owner          *Owner    `json:"owner"`
	Title          string    `json:"title"`
	Description    *string   `json:"description"`
	IsPinned       bool      `json:"isPinned"`
}

type GetAllWorkflowsResponse struct {
	Items      []*WorkflowResponse
	TotalCount int64
}
