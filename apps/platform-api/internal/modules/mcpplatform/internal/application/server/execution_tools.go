package server

import (
	"context"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	executionsContract "github.com/blocknextai/platform-api/internal/modules/executions/contract"
)

type executionOutput struct {
	ID            string     `json:"id" jsonschema:"execution id"`
	Status        string     `json:"status" jsonschema:"execution status"`
	ExecutionType string     `json:"executionType" jsonschema:"what kind of run produced this execution"`
	Context       string     `json:"context" jsonschema:"what was executed (workflow, library item, ...)"`
	ContextItemID string     `json:"contextItemId" jsonschema:"id of the executed item"`
	ErrorMessage  string     `json:"errorMessage,omitempty" jsonschema:"error message when the execution failed"`
	StartedAt     *time.Time `json:"startedAt,omitempty" jsonschema:"when the execution started"`
	CompletedAt   *time.Time `json:"completedAt,omitempty" jsonschema:"when the execution finished"`
}

type listExecutionsOutput struct {
	Total      int64             `json:"total" jsonschema:"total number of executions"`
	Executions []executionOutput `json:"executions" jsonschema:"the executions in this page"`
}

func (p *serverProvider) registerExecutionTools(server *mcpsdk.Server) {
	addTool(p, server, &mcpsdk.Tool{
		Name:        serverID + "_list_executions",
		Title:       "List executions",
		Description: "List the workflow executions of the organization the access token acts on, newest first.",
		Annotations: readOnly("List executions"),
	}, "play", platformReadScope, p.listExecutions)
}

func (p *serverProvider) listExecutions(ctx context.Context, req *mcpsdk.CallToolRequest, input listInput) (*mcpsdk.CallToolResult, listExecutionsOutput, error) {
	authenticated, err := requireCaller(req, platformReadScope)
	if err != nil {
		return nil, listExecutionsOutput{}, err
	}

	offset, limit := normalizePagination(input.Offset, input.Limit)

	executions, total, err := p.deps.TaskExecutionService.GetAllByOrganizationID(ctx, authenticated.OrganizationID, input.Search, offset, limit)
	if err != nil {
		return nil, listExecutionsOutput{}, err
	}

	outputs := make([]executionOutput, 0, len(executions))
	for _, execution := range executions {
		outputs = append(outputs, toExecutionOutput(execution))
	}

	return nil, listExecutionsOutput{
		Total:      total,
		Executions: outputs,
	}, nil
}

func toExecutionOutput(execution *executionsContract.TaskExecution) executionOutput {
	return executionOutput{
		ID:            execution.ID.String(),
		Status:        execution.Status,
		ExecutionType: string(execution.ExecutionType),
		Context:       string(execution.ExecutionContext),
		ContextItemID: execution.ContextItemID.String(),
		ErrorMessage:  stringValue(execution.ErrorMessage),
		StartedAt:     execution.StartedAt,
		CompletedAt:   execution.CompletedAt,
	}
}
