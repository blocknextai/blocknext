package taskrunner

import (
	"context"
	"log/slog"
	"sync"

	taskRunnerDomainTaskRunner "github.com/blocknextai/platform-api/internal/taskrunner/domain/taskrunner"
)

type WorkerPool struct {
	workers     int
	workerQueue chan chan func()
	jobQueue    chan func()
	quit        chan bool
	wg          sync.WaitGroup
}

func NewWorkerPool(workers int, jobQueueSize int) taskRunnerDomainTaskRunner.WorkerPool {
	wp := &WorkerPool{
		workers:     workers,
		workerQueue: make(chan chan func(), workers),
		jobQueue:    make(chan func(), jobQueueSize),
		quit:        make(chan bool),
	}

	return wp
}

func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < wp.workers; i++ {
		worker := NewWorker(i, wp.workerQueue, make(chan func()), make(chan bool))

		wp.wg.Add(1)
		go worker.Start(ctx, &wp.wg)
	}

	wp.wg.Add(1)
	go wp.dispatch(ctx, &wp.wg)
}

func (wp *WorkerPool) Submit(job func()) error {
	select {
	case wp.jobQueue <- job:
		return nil
	default:
		return ErrWorkerPoolFull
	}
}

func (wp *WorkerPool) SubmitWithContext(ctx context.Context, job func()) error {
	select {
	case wp.jobQueue <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.quit)
	wp.wg.Wait()
}

func (wp *WorkerPool) dispatch(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case job := <-wp.jobQueue:
			select {
			case workerChannel := <-wp.workerQueue:
				select {
				case workerChannel <- job:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		case <-wp.quit:
			return
		case <-ctx.Done():
			return
		}
	}
}

type Worker struct {
	ID          int
	WorkerQueue chan chan func()
	JobChannel  chan func()
	Quit        chan bool
}

func NewWorker(
	id int,
	workerQueue chan chan func(),
	jobChannel chan func(),
	quit chan bool,
) taskRunnerDomainTaskRunner.Worker {
	return &Worker{
		ID:          id,
		WorkerQueue: workerQueue,
		JobChannel:  jobChannel,
		Quit:        quit,
	}
}

func (w *Worker) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case w.WorkerQueue <- w.JobChannel:
			select {
			case job := <-w.JobChannel:
				func() {
					defer func() {
						if r := recover(); r != nil {
							slog.ErrorContext(ctx, "Worker panicked",
								"component", "worker_pool",
								"worker_id", w.ID,
								"panic", r)
						}
					}()
					job()
				}()
			case <-w.Quit:
				return
			case <-ctx.Done():
				return
			}
		case <-w.Quit:
			return
		case <-ctx.Done():
			return
		}
	}
}
