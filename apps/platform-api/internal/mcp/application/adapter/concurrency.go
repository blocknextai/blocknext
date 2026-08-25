package adapter

import (
	"context"
	"log/slog"
	"time"

	bnuuid "github.com/blocknextai/go-packages/uuid"
	"github.com/google/uuid"
)

func (a *adapter) acquireSlot(ctx context.Context, organizationID uuid.UUID) (func(), error) {
	holderID := bnuuid.NewV7()

	semaphore, err := a.semaphoreManager.AcquireSemaphore(ctx, organizationID, holderID, a.maxConcurrentExecutions)
	if err != nil {
		return nil, err
	}

	heartbeatCtx, stopHeartbeat := context.WithCancel(ctx)
	go a.runSemaphoreHeartbeat(heartbeatCtx, organizationID, holderID)

	releaseCtx := context.WithoutCancel(ctx)

	return func() {
		stopHeartbeat()
		if err := a.semaphoreManager.ReleaseSemaphore(releaseCtx, organizationID, holderID, semaphore); err != nil {
			slog.WarnContext(releaseCtx, "failed to release mcp semaphore",
				"component", "mcp_adapter",
				"organization_id", organizationID,
				"holder_id", holderID,
				"error", err)
		}
	}, nil
}

func (a *adapter) runSemaphoreHeartbeat(ctx context.Context, organizationID uuid.UUID, holderID uuid.UUID) {
	if a.heartbeatInterval <= 0 {
		return
	}

	ticker := time.NewTicker(a.heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := a.semaphoreManager.HeartbeatSemaphore(ctx, organizationID, holderID); err != nil {
				slog.WarnContext(ctx, "failed to heartbeat mcp semaphore",
					"component", "mcp_adapter",
					"organization_id", organizationID,
					"holder_id", holderID,
					"error", err)
			}
		}
	}
}
