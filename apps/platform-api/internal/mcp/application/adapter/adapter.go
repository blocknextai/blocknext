package adapter

import (
	"context"
	"errors"
	"time"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/json"
	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonDomainSemaphore "github.com/blocknextai/platform-api/internal/common/domain/semaphore"
	"github.com/blocknextai/platform-api/internal/mcp/application/credentialresolver"
	"github.com/blocknextai/platform-api/internal/mcp/application/history"
	nodeEngineApplicationExecutors "github.com/blocknextai/platform-api/internal/nodeengine/application/executors"
	nodeEngineDomainExecutors "github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	nodeEngineDomainMCP "github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/google/uuid"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	ErrExecutorNotFound          = apperror.NotFound("executor not found")
	ErrMissingOwner              = apperror.Internal("missing owner headers")
	ErrOrganizationOwnerRequired = apperror.Forbidden("mcp tools require an organization api key")
	ErrInvalidArguments          = apperror.Validation("invalid arguments")
	ErrToolPanicked              = apperror.Internal("tool execution panicked")
)

const (
	credentialsKey = "credentials"
	itemsKey       = "items"
)

type Adapter interface {
	Build(server nodeEngineDomainMCP.ServerManager) (*mcpsdk.Server, error)
}

type adapter struct {
	executorService         nodeEngineApplicationExecutors.ExecutorService
	credentialResolver      credentialresolver.CredentialResolver
	recorder                history.Recorder
	semaphoreManager        commonDomainSemaphore.SemaphoreManager
	maxConcurrentExecutions int64
	heartbeatInterval       time.Duration
	maxExecutionTime        time.Duration
}

func NewAdapter(
	executorService nodeEngineApplicationExecutors.ExecutorService,
	credentialResolver credentialresolver.CredentialResolver,
	recorder history.Recorder,
	semaphoreManager commonDomainSemaphore.SemaphoreManager,
	maxConcurrentExecutions int64,
	heartbeatInterval time.Duration,
	maxExecutionTime time.Duration,
) Adapter {
	return &adapter{
		executorService:         executorService,
		credentialResolver:      credentialResolver,
		recorder:                recorder,
		semaphoreManager:        semaphoreManager,
		maxConcurrentExecutions: maxConcurrentExecutions,
		heartbeatInterval:       heartbeatInterval,
		maxExecutionTime:        maxExecutionTime,
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
		}, a.buildHandler(nodeID, executor, credentialKeys))
	}

	return mcpServer, nil
}
func (a *adapter) buildHandler(
	nodeID string,
	executor nodeEngineDomainExecutors.ExecutorManager,
	credentialKeys []string,
) mcpsdk.ToolHandler {
	return func(ctx context.Context, req *mcpsdk.CallToolRequest) (*mcpsdk.CallToolResult, error) {
		ownerType, ownerID, err := extractOwner(req)
		if err != nil {
			return errorResult(err.Error()), nil
		}

		if ownerType != commonDomain.OwnerTypeOrganization {
			return errorResult(ErrOrganizationOwnerRequired.Error()), nil
		}

		input := map[string]any{}
		if len(req.Params.Arguments) > 0 {
			if err := json.Unmarshal(req.Params.Arguments, &input); err != nil {
				return errorResult(ErrInvalidArguments.WithCause(err).Error()), nil
			}
		}

		references, parameters := splitInput(input, credentialKeys)

		release, err := a.acquireSlot(ctx, ownerID)
		if err != nil {
			return errorResult(err.Error()), nil
		}
		defer release()

		callCtx := ctx
		if a.maxExecutionTime > 0 {
			var cancel context.CancelFunc
			callCtx, cancel = context.WithTimeout(ctx, a.maxExecutionTime)
			defer cancel()
		}

		recordCtx := context.WithoutCancel(ctx)
		call := history.ToolCall{
			OrganizationID: ownerID,
			APIKeyID:       extractAPIKeyID(req),
			ToolID:         nodeID,
			Parameters:     parameters,
			Credentials:    credentialReferences(references),
			StartedAt:      time.Now().UTC(),
		}

		var outputs []map[string]any

		defer func() {
			call.Err = err
			if recovered := recover(); recovered != nil {
				call.Err = ErrToolPanicked
				call.CompletedAt = time.Now().UTC()
				a.recorder.Record(recordCtx, call)
				panic(recovered)
			}

			call.Outputs = outputs
			call.CompletedAt = time.Now().UTC()
			a.recorder.Record(recordCtx, call)
		}()

		outputs, err = a.invoke(callCtx, ownerType, ownerID, executor, references, parameters)

		if err != nil {
			return errorResult(err.Error()), nil
		}

		return successResult(outputs), nil
	}
}

func (a *adapter) invoke(
	ctx context.Context,
	ownerType commonDomain.OwnerType,
	ownerID uuid.UUID,
	executor nodeEngineDomainExecutors.ExecutorManager,
	references map[string]any,
	parameters map[string]any,
) ([]map[string]any, error) {
	credentials, err := a.credentialResolver.Resolve(ctx, ownerType, ownerID, references)
	if err != nil {
		return nil, err
	}

	return executor.ExecuteWithContext(ctx, credentials, []map[string]any{parameters})
}
