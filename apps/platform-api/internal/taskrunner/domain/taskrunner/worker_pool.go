package taskrunner

import (
	"context"
	"sync"
)

type WorkerPool interface {
	Start(ctx context.Context)
	Submit(job func()) error
	SubmitWithContext(ctx context.Context, job func()) error
	Stop()
}

type Worker interface {
	Start(ctx context.Context, wg *sync.WaitGroup)
}
