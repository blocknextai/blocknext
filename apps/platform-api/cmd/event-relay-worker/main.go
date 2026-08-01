package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/blocknextai/platform-api/internal/bootstrap"
	"github.com/blocknextai/platform-api/internal/config"
)

func main() {
	cfg, err := config.LoadPlatformAPI()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	core, err := bootstrap.NewCore(cfg.Shared)
	if err != nil {
		slog.Error("failed to bootstrap core", "error", err)
		os.Exit(1)
	}

	app, err := bootstrap.NewPlatformAPI(core, cfg)
	if err != nil {
		slog.Error("failed to bootstrap module graph", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	core.EventBus.StartRelay(ctx)

	slog.Info("event-relay-worker started, draining transactional outbox")

	bootstrap.WaitForShutdown(
		cfg.TaskRunner.ShutdownTimeout,
		func() error { cancel(); return nil },
		app.Broadcaster.Close,
		core.Shutdown,
	)
}
