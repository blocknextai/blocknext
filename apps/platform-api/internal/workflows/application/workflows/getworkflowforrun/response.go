package getworkflowforrun

import (
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	nodeEngineDomainNodes "github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/google/uuid"
)

type GetWorkflowForRunResponse struct {
	ID                uuid.UUID                                       `json:"id"`
	OrganizationID    uuid.UUID                                       `json:"organizationId"`
	Title             string                                          `json:"title"`
	Description       *string                                         `json:"description"`
	Nodes             []RunNode                                       `json:"nodes"`
	CredentialSchemas []nodeEngineDomainCredentials.CredentialManager `json:"credentialSchemas"`
	NodeSchemas       []nodeEngineDomainNodes.NodeManager             `json:"nodeSchemas"`
}

type RunNode struct {
	ID     string `json:"id"`
	NodeID string `json:"nodeId"`
	Type   string `json:"type"`
}
