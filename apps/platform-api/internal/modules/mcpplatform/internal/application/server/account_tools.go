package server

import (
	"context"

	"github.com/google/uuid"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	organizationsContract "github.com/blocknextai/platform-api/internal/modules/organizations/contract"
)

type emptyInput struct{}

type organizationOutput struct {
	ID         string `json:"id" jsonschema:"organization id"`
	Title      string `json:"title" jsonschema:"organization title"`
	Role       string `json:"role" jsonschema:"the caller's role in the organization"`
	Alias      string `json:"alias,omitempty" jsonschema:"the caller's alias in the organization"`
	IsVerified bool   `json:"isVerified" jsonschema:"whether the organization is verified"`
	IsActive   bool   `json:"isActive" jsonschema:"whether the access token acts on this organization"`
}

type getMeOutput struct {
	UserID               string               `json:"userId" jsonschema:"the signed-in user id"`
	Role                 string               `json:"role" jsonschema:"the platform role of the user"`
	IsVerified           bool                 `json:"isVerified" jsonschema:"whether the user completed verification"`
	ActiveOrganizationID string               `json:"activeOrganizationId" jsonschema:"the organization the access token was issued for"`
	Organizations        []organizationOutput `json:"organizations" jsonschema:"organizations the user belongs to"`
}

func (p *serverProvider) registerAccountTools(server *mcpsdk.Server) {
	addTool(p, server, &mcpsdk.Tool{
		Name:        serverID + "_get_me",
		Title:       "Get my account",
		Description: "Return the signed-in user, the organization the access token acts on, and every organization the user belongs to.",
		Annotations: readOnly("Get my account"),
	}, "eye", platformReadScope, p.getMe)
}

func (p *serverProvider) getMe(ctx context.Context, req *mcpsdk.CallToolRequest, _ emptyInput) (*mcpsdk.CallToolResult, getMeOutput, error) {
	authenticated, err := requireCaller(req, platformReadScope)
	if err != nil {
		return nil, getMeOutput{}, err
	}

	user, err := p.deps.UserService.GetByID(ctx, authenticated.UserID)
	if err != nil {
		return nil, getMeOutput{}, err
	}

	memberships, err := p.deps.OrganizationUserService.GetAllByUserID(ctx, authenticated.UserID)
	if err != nil {
		return nil, getMeOutput{}, err
	}

	organizationIDs := make([]uuid.UUID, 0, len(memberships))
	for _, membership := range memberships {
		organizationIDs = append(organizationIDs, membership.OrganizationID)
	}

	organizations := make([]*organizationsContract.Organization, 0, len(organizationIDs))
	if len(organizationIDs) > 0 {
		organizations, err = p.deps.OrganizationService.GetAllByIDs(ctx, organizationIDs)
		if err != nil {
			return nil, getMeOutput{}, err
		}
	}

	organizationsByID := make(map[uuid.UUID]*organizationsContract.Organization, len(organizations))
	for _, organization := range organizations {
		organizationsByID[organization.ID] = organization
	}

	outputs := make([]organizationOutput, 0, len(memberships))
	for _, membership := range memberships {
		organization, ok := organizationsByID[membership.OrganizationID]
		if !ok {
			continue
		}

		outputs = append(outputs, organizationOutput{
			ID:         organization.ID.String(),
			Title:      organization.Title,
			Role:       membership.Role,
			Alias:      membership.Alias,
			IsVerified: organization.IsVerified,
			IsActive:   organization.ID == authenticated.OrganizationID,
		})
	}

	return nil, getMeOutput{
		UserID:               user.ID.String(),
		Role:                 user.Role,
		IsVerified:           user.IsVerified,
		ActiveOrganizationID: authenticated.OrganizationID.String(),
		Organizations:        outputs,
	}, nil
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
