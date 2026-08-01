package bootstrap

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func WaitForShutdown(timeout time.Duration, shutdownFns ...func() error) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	slog.Info("shutting down", "timeout", timeout)

	done := make(chan struct{})
	go func() {
		for _, fn := range shutdownFns {
			if fn == nil {
				continue
			}
			if err := fn(); err != nil {
				slog.Error("shutdown step failed", "error", err)
			}
		}
		close(done)
	}()

	select {
	case <-done:
		slog.Info("shutdown complete")
	case <-time.After(timeout):
		slog.Warn("shutdown timeout exceeded; forcing exit")
	}
}
