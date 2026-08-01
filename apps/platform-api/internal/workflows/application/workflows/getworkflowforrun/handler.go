package getworkflowforrun

import (
	"context"

	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
	nodeEngineApplicationNodes "github.com/blocknextai/platform-api/internal/nodeengine/application/nodes"
	workflowsDomainWorkflows "github.com/blocknextai/platform-api/internal/workflows/domain/workflows"
)

type Handler struct {
	workflowRepository workflowsDomainWorkflows.WorkflowRepository
	credentialService  nodeEngineApplicationCredentials.CredentialService
	nodeService        nodeEngineApplicationNodes.NodeService
}

func New(
	workflowRepository workflowsDomainWorkflows.WorkflowRepository,
	credentialService nodeEngineApplicationCredentials.CredentialService,
	nodeService nodeEngineApplicationNodes.NodeService,
) *Handler {
	return &Handler{
		workflowRepository: workflowRepository,
		credentialService:  credentialService,
		nodeService:        nodeService,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetWorkflowForRunQuery) (*GetWorkflowForRunResponse, error) {
	workflow, err := h.workflowRepository.GetByOrganizationIDAndID(ctx, request.OrganizationID, request.WorkflowID)
	if err != nil {
		return nil, err
	}

	nodes := workflow.Nodes

	nodeIDs := make([]string, 0)
	for _, node := range nodes {
		nodeIDs = append(nodeIDs, node.NodeID)
	}

	credentialSchemas := h.credentialService.GetCredentialSchemasByNodeIDs(nodeIDs)
	nodeSchemas := h.nodeService.GetNodeSchemasByIDs(nodeIDs)

	return MapWorkflowToResponse(workflow, nodes, credentialSchemas, nodeSchemas), nil
}
