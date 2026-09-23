package server

import (
	"context"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	executionsContract "github.com/blocknextai/platform-api/internal/modules/executions/contract"
)

type toolInvocationOutput struct {
	ID           string    `json:"id" jsonschema:"tool invocation id"`
	Source       string    `json:"source" jsonschema:"which surface made the call (mcp)"`
	ToolID       string    `json:"toolId" jsonschema:"the tool that was called"`
	Status       string    `json:"status" jsonschema:"success or failed"`
	ErrorMessage string    `json:"errorMessage,omitempty" jsonschema:"error message when the call failed"`
	StartedAt    time.Time `json:"startedAt" jsonschema:"when the call started"`
	CompletedAt  time.Time `json:"completedAt" jsonschema:"when the call finished"`
}

type listToolInvocationsOutput struct {
	Total           int64                  `json:"total" jsonschema:"total number of tool invocations"`
	ToolInvocations []toolInvocationOutput `json:"toolInvocations" jsonschema:"the tool invocations in this page"`
}

func (p *serverProvider) registerToolInvocationTools(server *mcpsdk.Server) {
	addTool(p, server, &mcpsdk.Tool{
		Name:        serverID + "_list_tool_invocations",
		Title:       "List MCP tool calls",
		Description: "List the tool calls made through the MCP servers on behalf of the organization the access token acts on, newest first. Parameters and outputs are not included.",
		Annotations: readOnly("List MCP tool calls"),
	}, "record", platformReadScope, p.listToolInvocations)
}

func (p *serverProvider) listToolInvocations(ctx context.Context, req *mcpsdk.CallToolRequest, input listInput) (*mcpsdk.CallToolResult, listToolInvocationsOutput, error) {
	authenticated, err := requireCaller(req, platformReadScope)
	if err != nil {
		return nil, listToolInvocationsOutput{}, err
	}

	offset, limit := normalizePagination(input.Offset, input.Limit)

	toolInvocations, total, err := p.deps.ToolInvocationService.GetAllByOrganizationID(ctx, authenticated.OrganizationID, input.Search, offset, limit)
	if err != nil {
		return nil, listToolInvocationsOutput{}, err
	}

	outputs := make([]toolInvocationOutput, 0, len(toolInvocations))
	for _, toolInvocation := range toolInvocations {
		outputs = append(outputs, toToolInvocationOutput(toolInvocation))
	}

	return nil, listToolInvocationsOutput{
		Total:           total,
		ToolInvocations: outputs,
	}, nil
}

func toToolInvocationOutput(toolInvocation *executionsContract.ToolInvocation) toolInvocationOutput {
	return toolInvocationOutput{
		ID:           toolInvocation.ID.String(),
		Source:       string(toolInvocation.Source),
		ToolID:       toolInvocation.ToolID,
		Status:       string(toolInvocation.Status),
		ErrorMessage: stringValue(toolInvocation.ErrorMessage),
		StartedAt:    toolInvocation.StartedAt,
		CompletedAt:  toolInvocation.CompletedAt,
	}
}
