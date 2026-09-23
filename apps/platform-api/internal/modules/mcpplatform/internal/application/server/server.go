package server

import (
	"context"
	"log/slog"
	"time"

	gjs "github.com/google/jsonschema-go/jsonschema"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/blocknextai/go-packages/json"

	accountContract "github.com/blocknextai/platform-api/internal/modules/account/contract"
	credentialsContract "github.com/blocknextai/platform-api/internal/modules/credentials/contract"
	executionsContract "github.com/blocknextai/platform-api/internal/modules/executions/contract"
	mcpContract "github.com/blocknextai/platform-api/internal/modules/mcp/contract"
	organizationsContract "github.com/blocknextai/platform-api/internal/modules/organizations/contract"
	triggersContract "github.com/blocknextai/platform-api/internal/modules/triggers/contract"
	workflowsContract "github.com/blocknextai/platform-api/internal/modules/workflows/contract"
)

const (
	serverID          = "platform"
	serverName        = "Platform"
	serverBrand       = "platform"
	serverDescription = "Read and manage your account: profile, organizations, workflows, executions, MCP tool calls and credentials."

	defaultPageSize = 20
	maxPageSize     = 100

	instructions = "Manage the signed-in account: read the profile and organizations, " +
		"list and inspect workflows, triggers, executions, MCP tool calls and credentials, " +
		"and turn triggers on or off. Every tool acts on the organization the access token was issued for."
)

type Dependencies struct {
	Version string

	UserService             accountContract.UserService
	OrganizationService     organizationsContract.OrganizationService
	OrganizationUserService organizationsContract.OrganizationUserService
	WorkflowService         workflowsContract.WorkflowService
	TriggerService          triggersContract.TriggerService
	TaskExecutionService    executionsContract.TaskExecutionService
	ToolInvocationService   executionsContract.ToolInvocationService
	CredentialService       credentialsContract.CredentialService
}

type serverProvider struct {
	deps   Dependencies
	server *mcpsdk.Server
	tools  []mcpContract.Tool
}

func NewServerProvider(deps Dependencies) mcpContract.ServerProvider {
	provider := &serverProvider{
		deps: deps,
	}
	provider.server = provider.build()

	return provider
}

func (p *serverProvider) build() *mcpsdk.Server {
	server := mcpsdk.NewServer(&mcpsdk.Implementation{
		Name:    serverID,
		Title:   serverName,
		Version: p.deps.Version,
	}, &mcpsdk.ServerOptions{
		Instructions: instructions,
	})

	p.registerAccountTools(server)
	p.registerWorkflowTools(server)
	p.registerTriggerTools(server)
	p.registerExecutionTools(server)
	p.registerToolInvocationTools(server)
	p.registerCredentialTools(server)

	return server
}

func (p *serverProvider) GetAllServers() ([]*mcpContract.Server, error) {
	return []*mcpContract.Server{
		{
			ID:           serverID,
			Name:         serverName,
			Description:  serverDescription,
			Icon:         mcpContract.Icon{Brand: serverBrand},
			Version:      p.deps.Version,
			Instructions: instructions,
			AuthMethods: []mcpContract.AuthMethod{
				mcpContract.AuthMethodOAuth,
			},
			Scopes: []string{
				platformReadScope,
				platformWriteScope,
			},
			Tools:     p.tools,
			MCPServer: p.server,
		},
	}, nil
}

func addTool[In, Out any](
	provider *serverProvider,
	server *mcpsdk.Server,
	tool *mcpsdk.Tool,
	glyph string,
	scope string,
	handler mcpsdk.ToolHandlerFor[In, Out],
) {
	mcpsdk.AddTool(server, tool, recorded(provider, tool.Name, handler))

	provider.tools = append(provider.tools, mcpContract.Tool{
		ID:           tool.Name,
		Version:      provider.deps.Version,
		Name:         tool.Title,
		Description:  tool.Description,
		Icon:         mcpContract.Icon{Brand: serverBrand, Glyph: glyph},
		InputSchema:  inferSchema[In](tool.Name, "input"),
		OutputSchema: inferSchema[Out](tool.Name, "output"),
		Scopes:       []string{scope},
	})
}

func recorded[In, Out any](
	provider *serverProvider,
	toolID string,
	handler mcpsdk.ToolHandlerFor[In, Out],
) mcpsdk.ToolHandlerFor[In, Out] {
	return func(ctx context.Context, req *mcpsdk.CallToolRequest, input In) (*mcpsdk.CallToolResult, Out, error) {
		startedAt := time.Now().UTC()

		result, output, err := handler(ctx, req, input)

		provider.record(context.WithoutCancel(ctx), req, toolID, input, output, err, startedAt)

		return result, output, err
	}
}

func (p *serverProvider) record(
	ctx context.Context,
	req *mcpsdk.CallToolRequest,
	toolID string,
	input any,
	output any,
	callErr error,
	startedAt time.Time,
) {
	authenticated, err := callerFromRequest(req)
	if err != nil {
		return
	}

	status := executionsContract.StatusSuccess
	var errorMessage *string
	var outputs []map[string]any
	if callErr != nil {
		status = executionsContract.StatusFailed
		errorMessage = new(callErr.Error())
	} else {
		outputs = []map[string]any{toMap(output)}
	}

	if _, err := p.deps.ToolInvocationService.Record(
		ctx,
		authenticated.OrganizationID,
		nil,
		executionsContract.SourcePlatform,
		toolID,
		status,
		toMap(input),
		nil,
		outputs,
		errorMessage,
		startedAt,
		time.Now().UTC(),
	); err != nil {
		slog.WarnContext(ctx, "failed to record platform tool call",
			"component", "mcpplatform",
			"organization_id", authenticated.OrganizationID,
			"tool_id", toolID,
			"error", err)
	}
}

func toMap(value any) map[string]any {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil
	}

	out := map[string]any{}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil
	}

	return out
}

func inferSchema[T any](toolName string, kind string) *gjs.Schema {
	schema, err := gjs.For[T](nil)
	if err != nil {
		slog.Warn("failed to infer platform mcp tool schema",
			"component", "mcpplatform",
			"tool", toolName,
			"schema", kind,
			"error", err.Error(),
		)
		return nil
	}

	return schema
}

func readOnly(title string) *mcpsdk.ToolAnnotations {
	return &mcpsdk.ToolAnnotations{
		Title:        title,
		ReadOnlyHint: true,
	}
}

func mutating(title string) *mcpsdk.ToolAnnotations {
	destructive := false

	return &mcpsdk.ToolAnnotations{
		Title:           title,
		DestructiveHint: &destructive,
		IdempotentHint:  true,
	}
}
