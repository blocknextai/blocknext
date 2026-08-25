package realtime

import (
	"context"
	executionsDomainToolInvocations "github.com/blocknextai/platform-api/internal/executions/domain/toolinvocations"

	taskRunnerDomainNode "github.com/blocknextai/platform-api/internal/taskrunner/domain/node"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	"github.com/google/uuid"
)

type Broadcaster interface {
	Ping(ctx context.Context) error
	PublishTaskEvent(ctx context.Context, event *taskRunnerDomainTask.TaskEvent) error
	PublishNodeEvent(ctx context.Context, event *taskRunnerDomainNode.NodeEvent) error
	PublishToolInvocationEvent(ctx context.Context, event *executionsDomainToolInvocations.ToolInvocationEvent) error
	Subscribe(ctx context.Context, organizationID uuid.UUID) (<-chan string, error)
	Close() error
}
