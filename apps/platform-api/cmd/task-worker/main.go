package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/blocknextai/platform-api/internal/bootstrap"
	"github.com/blocknextai/platform-api/internal/config"
)

func main() {
	cfg, err := config.LoadTaskWorker()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	if cfg.TaskRunner.Mode != config.TaskRunnerModeQueue {
		slog.Info("task-worker is disabled", "mode", string(cfg.TaskRunner.Mode))
		return
	}

	core, err := bootstrap.NewCore(cfg.Shared)
	if err != nil {
		slog.Error("failed to bootstrap core", "error", err)
		os.Exit(1)
	}

	stack, err := bootstrap.NewTaskWorker(core, cfg)
	if err != nil {
		slog.Error("failed to bootstrap task-worker stack", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := stack.TaskRunnerModule.StartAsWorker(ctx); err != nil {
		slog.Error("failed to start task worker", "error", err)
		os.Exit(1)
	}

	slog.Info("task-worker started, waiting for tasks from queue")

	bootstrap.WaitForShutdown(
		cfg.TaskRunner.ShutdownTimeout,
		func() error { cancel(); return nil },
		stack.Shutdown,
		core.Shutdown,
	)
}
