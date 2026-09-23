package taskrunner

import (
	"context"
)

type LeaderRunner interface {
	Ping(ctx context.Context) error
	Run(ctx context.Context, leaderFunc func(leaderCtx context.Context))
}
