package server

import (
	"context"
	"time"

	"github.com/google/uuid"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	workflowsContract "github.com/blocknextai/platform-api/internal/modules/workflows/contract"
)

type listInput struct {
	Search string `json:"search,omitempty" jsonschema:"search text"`
	Offset int    `json:"offset,omitempty" jsonschema:"number of records to skip"`
	Limit  int    `json:"limit,omitempty" jsonschema:"maximum number of records to return (default 20, max 100)"`
}

type workflowOutput struct {
	ID          string    `json:"id" jsonschema:"workflow id"`
	Title       string    `json:"title" jsonschema:"workflow title"`
	Description string    `json:"description,omitempty" jsonschema:"workflow description"`
	IsPinned    bool      `json:"isPinned" jsonschema:"whether the workflow is pinned in the UI"`
	NodeCount   int       `json:"nodeCount" jsonschema:"number of nodes on the canvas"`
	UpdatedAt   time.Time `json:"updatedAt" jsonschema:"last update time"`
}

type listWorkflowsOutput struct {
	Total     int64            `json:"total" jsonschema:"total number of workflows"`
	Workflows []workflowOutput `json:"workflows" jsonschema:"the workflows in this page"`
}

type getWorkflowInput struct {
	WorkflowID string `json:"workflowId" jsonschema:"the workflow id"`
}

type workflowNodeOutput struct {
	ID     string `json:"id" jsonschema:"node id on the canvas"`
	Type   string `json:"type" jsonschema:"node kind on the canvas"`
	NodeID string `json:"nodeId" jsonschema:"catalog node the canvas node runs"`
	Title  string `json:"title,omitempty" jsonschema:"node title"`
}

type getWorkflowOutput struct {
	Workflow workflowOutput       `json:"workflow" jsonschema:"the workflow"`
	Nodes    []workflowNodeOutput `json:"nodes" jsonschema:"the nodes of the workflow"`
}

func (p *serverProvider) registerWorkflowTools(server *mcpsdk.Server) {
	addTool(p, server, &mcpsdk.Tool{
		Name:        serverID + "_list_workflows",
		Title:       "List workflows",
		Description: "List the workflows of the organization the access token acts on.",
		Annotations: readOnly("List workflows"),
	}, platformReadScope, p.listWorkflows)

	addTool(p, server, &mcpsdk.Tool{
		Name:        serverID + "_get_workflow",
		Title:       "Get workflow",
		Description: "Return one workflow with the nodes on its canvas.",
		Annotations: readOnly("Get workflow"),
	}, platformReadScope, p.getWorkflow)
}

func (p *serverProvider) listWorkflows(ctx context.Context, req *mcpsdk.CallToolRequest, input listInput) (*mcpsdk.CallToolResult, listWorkflowsOutput, error) {
	authenticated, err := requireCaller(req, platformReadScope)
	if err != nil {
		return nil, listWorkflowsOutput{}, err
	}

	offset, limit := normalizePagination(input.Offset, input.Limit)

	workflows, total, err := p.deps.WorkflowService.GetAllByOrganizationID(ctx, authenticated.OrganizationID, input.Search, offset, limit)
	if err != nil {
		return nil, listWorkflowsOutput{}, err
	}

	outputs := make([]workflowOutput, 0, len(workflows))
	for _, workflow := range workflows {
		outputs = append(outputs, toWorkflowOutput(workflow))
	}

	return nil, listWorkflowsOutput{
		Total:     total,
		Workflows: outputs,
	}, nil
}

func (p *serverProvider) getWorkflow(ctx context.Context, req *mcpsdk.CallToolRequest, input getWorkflowInput) (*mcpsdk.CallToolResult, getWorkflowOutput, error) {
	authenticated, err := requireCaller(req, platformReadScope)
	if err != nil {
		return nil, getWorkflowOutput{}, err
	}

	workflowID, err := uuid.Parse(input.WorkflowID)
	if err != nil {
		return nil, getWorkflowOutput{}, ErrInvalidWorkflowID
	}

	workflow, err := p.deps.WorkflowService.GetWorkflow(ctx, authenticated.OrganizationID, workflowID)
	if err != nil {
		return nil, getWorkflowOutput{}, err
	}

	nodes := make([]workflowNodeOutput, 0, len(workflow.Nodes))
	for _, node := range workflow.Nodes {
		nodes = append(nodes, workflowNodeOutput{
			ID:     node.ID,
			Type:   node.Type,
			NodeID: node.NodeID,
			Title:  node.Title,
		})
	}

	return nil, getWorkflowOutput{
		Workflow: toWorkflowOutput(workflow),
		Nodes:    nodes,
	}, nil
}

func toWorkflowOutput(workflow *workflowsContract.Workflow) workflowOutput {
	return workflowOutput{
		ID:          workflow.ID.String(),
		Title:       workflow.Title,
		Description: stringValue(workflow.Description),
		IsPinned:    workflow.IsPinned,
		NodeCount:   len(workflow.Nodes),
		UpdatedAt:   workflow.UpdatedAt,
	}
}
