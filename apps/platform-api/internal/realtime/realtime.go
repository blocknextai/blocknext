package realtime

import (
	"context"

	"github.com/google/uuid"

	"github.com/blocknextai/platform-api/internal/realtime/events"
)

type Broadcaster interface {
	Ping(ctx context.Context) error
	PublishTaskEvent(ctx context.Context, event *events.TaskEvent) error
	PublishNodeEvent(ctx context.Context, event *events.NodeEvent) error
	PublishToolInvocationEvent(ctx context.Context, event *events.ToolInvocationEvent) error
	Subscribe(ctx context.Context, organizationID uuid.UUID) (<-chan string, error)
	Close() error
}
