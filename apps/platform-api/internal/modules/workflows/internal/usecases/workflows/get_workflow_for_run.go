package workflows

import (
	"context"

	"github.com/google/uuid"

	"github.com/blocknextai/go-packages/dag"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	workflowsDomainWorkflows "github.com/blocknextai/platform-api/internal/modules/workflows/internal/domain/workflows"
)

type GetWorkflowForRunQuery struct {
	OrganizationID uuid.UUID
	WorkflowID     uuid.UUID
}

type GetWorkflowForRunResponse struct {
	ID                uuid.UUID                              `json:"id"`
	OrganizationID    uuid.UUID                              `json:"organizationId"`
	Title             string                                 `json:"title"`
	Description       *string                                `json:"description"`
	Nodes             []RunNode                              `json:"nodes"`
	CredentialSchemas []nodeengineContract.CredentialManager `json:"credentialSchemas"`
	NodeSchemas       []nodeengineContract.NodeManager       `json:"nodeSchemas"`
}

type RunNode struct {
	ID     string `json:"id"`
	NodeID string `json:"nodeId"`
	Type   string `json:"type"`
}

func (s *Service) GetWorkflowForRun(ctx context.Context, request *GetWorkflowForRunQuery) (*GetWorkflowForRunResponse, error) {
	workflow, err := s.workflowRepository.GetByOrganizationIDAndID(ctx, request.OrganizationID, request.WorkflowID)
	if err != nil {
		return nil, err
	}

	nodes := workflow.Nodes

	nodeIDs := make([]string, 0)
	for _, node := range nodes {
		nodeIDs = append(nodeIDs, node.NodeID)
	}

	credentialSchemas := s.credentialService.GetCredentialSchemasByNodeIDs(nodeIDs)
	nodeSchemas := s.nodeService.GetNodeSchemasByIDs(nodeIDs)

	return GetWorkflowForRunMapWorkflowToResponse(workflow, nodes, credentialSchemas, nodeSchemas), nil
}

func GetWorkflowForRunMapWorkflowToResponse(
	workflow *workflowsDomainWorkflows.Workflow,
	nodes []dag.Node,
	credentialSchemas []nodeengineContract.CredentialManager,
	nodeSchemas []nodeengineContract.NodeManager,
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
