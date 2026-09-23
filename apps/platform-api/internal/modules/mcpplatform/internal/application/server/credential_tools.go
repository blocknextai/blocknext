package server

import (
	"context"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	credentialsContract "github.com/blocknextai/platform-api/internal/modules/credentials/contract"
)

type credentialOutput struct {
	ID         string    `json:"id" jsonschema:"credential id, usable as a credential reference"`
	Key        string    `json:"key" jsonschema:"the credential type from the node catalog"`
	Title      string    `json:"title" jsonschema:"credential title"`
	SourceType string    `json:"sourceType" jsonschema:"how the credential was provided"`
	UpdatedAt  time.Time `json:"updatedAt" jsonschema:"last update time"`
}

type listCredentialsOutput struct {
	Total       int64              `json:"total" jsonschema:"total number of credentials"`
	Credentials []credentialOutput `json:"credentials" jsonschema:"the credentials in this page"`
}

func (p *serverProvider) registerCredentialTools(server *mcpsdk.Server) {
	addTool(p, server, &mcpsdk.Tool{
		Name:        serverID + "_list_credentials",
		Title:       "List credentials",
		Description: "List the credentials connected to the organization the access token acts on. Secrets are never returned.",
		Annotations: readOnly("List credentials"),
	}, platformReadScope, p.listCredentials)
}

func (p *serverProvider) listCredentials(ctx context.Context, req *mcpsdk.CallToolRequest, input listInput) (*mcpsdk.CallToolResult, listCredentialsOutput, error) {
	authenticated, err := requireCaller(req, platformReadScope)
	if err != nil {
		return nil, listCredentialsOutput{}, err
	}

	offset, limit := normalizePagination(input.Offset, input.Limit)

	credentials, total, err := p.deps.CredentialService.GetAllByOwner(
		ctx,
		commonDomain.OwnerTypeOrganization,
		authenticated.OrganizationID,
		input.Search,
		offset,
		limit,
	)
	if err != nil {
		return nil, listCredentialsOutput{}, err
	}

	outputs := make([]credentialOutput, 0, len(credentials))
	for _, credential := range credentials {
		outputs = append(outputs, toCredentialOutput(credential))
	}

	return nil, listCredentialsOutput{
		Total:       total,
		Credentials: outputs,
	}, nil
}

func toCredentialOutput(credential *credentialsContract.Credential) credentialOutput {
	return credentialOutput{
		ID:         credential.ID.String(),
		Key:        credential.Key,
		Title:      credential.Title,
		SourceType: string(credential.SourceType),
		UpdatedAt:  credential.UpdatedAt,
	}
}
