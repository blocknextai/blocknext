package getworkflowforrun

import (
	"github.com/blocknextai/go-packages/dag"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	nodeEngineDomainNodes "github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	workflowsDomainWorkflows "github.com/blocknextai/platform-api/internal/workflows/domain/workflows"
)

func MapWorkflowToResponse(
	workflow *workflowsDomainWorkflows.Workflow,
	nodes []dag.Node,
	credentialSchemas []nodeEngineDomainCredentials.CredentialManager,
	nodeSchemas []nodeEngineDomainNodes.NodeManager,
) *GetWorkflowForRunResponse {
	runNodes := MapNodesToRunNodes(nodes)

	return &GetWorkflowForRunResponse{
		ID:                workflow.ID,
		OrganizationID:    workflow.OrganizationID,
		Title:             workflow.Title,
		Description:       workflow.Description,
		Nodes:             runNodes,
		CredentialSchemas: credentialSchemas,
		NodeSchemas:       nodeSchemas,
	}
}

func MapNodesToRunNodes(nodes []dag.Node) []RunNode {
	runNodes := make([]RunNode, 0, len(nodes))
	for _, node := range nodes {
		runNodes = append(runNodes, RunNode{
			ID:     node.ID,
			NodeID: node.NodeID,
			Type:   node.Type,
		})
	}
	return runNodes
}
