package adapter

import (
	"context"
	"errors"
	"maps"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/json"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	"github.com/blocknextai/platform-api/internal/mcp/application/credentialresolver"
	nodeEngineApplicationExecutors "github.com/blocknextai/platform-api/internal/nodeengine/application/executors"
	nodeEngineDomainExecutors "github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	nodeEngineDomainMCP "github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	nodeEngineDomainNodes "github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
	"github.com/google/uuid"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	ErrExecutorNotFound = apperror.NotFound("executor not found")
	ErrMissingOwner     = apperror.Internal("missing owner headers")
)

const (
	credentialsKey = "credentials"
	itemsKey       = "items"
)

type Adapter interface {
	Build(server nodeEngineDomainMCP.ServerManager) (*mcpsdk.Server, error)
}

type adapter struct {
	executorService    nodeEngineApplicationExecutors.ExecutorService
	credentialResolver credentialresolver.CredentialResolver
}

func NewAdapter(
	executorService nodeEngineApplicationExecutors.ExecutorService,
	credentialResolver credentialresolver.CredentialResolver,
) Adapter {
	return &adapter{
		executorService:    executorService,
		credentialResolver: credentialResolver,
	}
}

func (a *adapter) Build(server nodeEngineDomainMCP.ServerManager) (*mcpsdk.Server, error) {
	mcpServer := mcpsdk.NewServer(&mcpsdk.Implementation{
		Name:    server.GetID(),
		Title:   server.GetName(),
		Version: server.GetVersion(),
	}, &mcpsdk.ServerOptions{
		Instructions: server.GetInstructions(),
	})

	for _, node := range server.GetTools() {
		nodeID := node.GetID()
		executor, ok := a.executorService.GetByID(nodeID)
		if !ok {
			return nil, ErrExecutorNotFound.WithCause(errors.New(nodeID))
		}

		credentialKeys := node.GetSupportedCredentials()

		mcpServer.AddTool(&mcpsdk.Tool{
			Name:         nodeID,
			Title:        node.GetName(),
			Description:  node.GetDescription(),
			Annotations:  mcpAnnotations(node.GetAnnotations()),
			InputSchema:  augmentInputSchema(node.GetInputSchema(), credentialKeys),
			OutputSchema: wrapOutputSchema(node.GetOutputSchema()),
		}, a.buildHandler(executor, credentialKeys))
	}

	return mcpServer, nil
}

func mcpAnnotations(a nodeEngineDomainNodes.NodeAnnotations) *mcpsdk.ToolAnnotations {
	return &mcpsdk.ToolAnnotations{
		ReadOnlyHint:    a.ReadOnly,
		DestructiveHint: a.Destructive,
		IdempotentHint:  a.Idempotent,
		OpenWorldHint:   a.OpenWorld,
	}
}

func augmentInputSchema(original *gjs.Schema, credentialKeys []string) *gjs.Schema {
	if len(credentialKeys) == 0 {
		return original
	}

	base := &gjs.Schema{Type: "object", Properties: map[string]*gjs.Schema{}}
	if original != nil {
		base = new(*original)
	}

	properties := make(map[string]*gjs.Schema, len(base.Properties)+1)
	maps.Copy(properties, base.Properties)

	credentialProps := make(map[string]*gjs.Schema, len(credentialKeys))
	for _, key := range credentialKeys {
		credentialProps[key] = &gjs.Schema{
			Type:        "string",
			Description: "Reference to a saved credential. Format: credential:organization:<uuid> or credential:user:<uuid>.",
			Pattern:     "^credential:(organization|user):[0-9a-fA-F-]+$",
		}
	}

	properties[credentialsKey] = &gjs.Schema{
		Type:        "object",
		Title:       "Credentials",
		Description: "References to saved credentials required by this tool.",
		Properties:  credentialProps,
		Required:    append([]string{}, credentialKeys...),
	}

	required := append([]string{}, base.Required...)
	required = append(required, credentialsKey)

	base.Properties = properties
	base.Required = required
	return base
}

func wrapOutputSchema(original *gjs.Schema) *gjs.Schema {
	if original == nil {
		return nil
	}
	return &gjs.Schema{
		Type: "object",
		Properties: map[string]*gjs.Schema{
			itemsKey: original,
		},
		Required: []string{
			itemsKey,
		},
	}
}

func (a *adapter) buildHandler(
	executor nodeEngineDomainExecutors.ExecutorManager,
	credentialKeys []string,
) mcpsdk.ToolHandler {
	return func(ctx context.Context, req *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
		ownerType, ownerID, err := extractOwner(req)
		if err != nil {
			return errorResult(err.Error()), nil
		}

		input := map[string]any{}
		if len(req.Params.Arguments) > 0 {
			if err := json.Unmarshal(req.Params.Arguments, &input); err != nil {
				return errorResult("invalid arguments: " + err.Error()), nil
			}
		}

		references, data := splitInput(input, credentialKeys)

		credentials, err := a.credentialResolver.Resolve(ctx, ownerType, ownerID, references)
		if err != nil {
			return errorResult(err.Error()), nil
		}

		outputs, err := executor.ExecuteWithContext(ctx, credentials, []map[string]any{data})
		if err != nil {
			return errorResult(err.Error()), nil
		}

		return successResult(outputs), nil
	}
}

func extractOwner(req *mcpsdk.CallToolRequest) (commonDomain.OwnerType, uuid.UUID, error) {
	if req.Extra == nil {
		return "", uuid.Nil, ErrMissingOwner
	}

	ownerType := commonDomain.OwnerType(req.Extra.Header.Get(commonAuth.OwnerTypeHeader))
	if ownerType == "" {
		return "", uuid.Nil, ErrMissingOwner
	}

	ownerID, err := uuid.Parse(req.Extra.Header.Get(commonAuth.OwnerIDHeader))
	if err != nil {
		return "", uuid.Nil, ErrMissingOwner.WithCause(err)
	}

	return ownerType, ownerID, nil
}

func splitInput(input map[string]any, credentialKeys []string) (map[string]any, map[string]any) {
	if len(credentialKeys) == 0 {
		return map[string]any{}, input
	}

	references := map[string]any{}
	if raw, ok := input[credentialsKey].(map[string]any); ok {
		references = raw
	}

	data := make(map[string]any, len(input))
	for k, v := range input {
		if k == credentialsKey {
			continue
		}
		data[k] = v
	}

	return references, data
}

func errorResult(message string) *mcpsdk.CallToolResult {
	return &mcpsdk.CallToolResult{
		IsError: true,
		Content: []mcpsdk.Content{
			&mcpsdk.TextContent{Text: message},
		},
	}
}

func successResult(outputs []map[string]any) *mcpsdk.CallToolResult {
	content := make([]mcpsdk.Content, 0, len(outputs))
	for _, output := range outputs {
		text, err := json.Marshal(output)
		if err != nil {
			return errorResult("failed to marshal output: " + err.Error())
		}
		content = append(content, &mcpsdk.TextContent{Text: string(text)})
	}

	return &mcpsdk.CallToolResult{
		StructuredContent: map[string]any{
			itemsKey: outputs,
		},
		Content: content,
	}
}
