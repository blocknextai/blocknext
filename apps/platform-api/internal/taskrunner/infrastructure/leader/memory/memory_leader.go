package memory

import (
	"context"
)

type MemoryLeader struct{}

func New() *MemoryLeader {
	return &MemoryLeader{}
}

func (l *MemoryLeader) Ping(_ context.Context) error {
	return nil
}

func (l *MemoryLeader) Run(ctx context.Context, leaderFunc func(leaderCtx context.Context)) {
	leaderFunc(ctx)
}
