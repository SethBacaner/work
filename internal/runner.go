package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/sethbacaner/work/internal/pool"
	"github.com/sethbacaner/work/internal/queue"
)

type Runner interface {
	Run(ctx context.Context, pool pool.Pool, queue queue.Queue, registry *Registry)
}

func NewRunner() Runner {
	return &RunnerImpl{}
}

type RunnerImpl struct{}

func (ri *RunnerImpl) Run(ctx context.Context, pool pool.Pool, queue queue.Queue, registry *Registry) {
	for {
		select {
		case <-ctx.Done():
			return
		// TODO: receive message on channel indicating there's an available worker
		case <-time.After(1 * time.Second):
			if pool.AvailableWorkers() <= 0 {
				continue
			}

			serializedTask, err := queue.Dequeue(ctx)

			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			if serializedTask == "" {
				continue
			}

			pool.Execute(func() {
				DoTask(ctx, serializedTask, *registry)
			})
		}
	}
}

// TODO: task might need to be reenqueued here.
func DoTask(ctx context.Context, serializedTask string, registry Registry) {

	args, definition, err := Deserialize(serializedTask, registry)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = definition.TaskFn(ctx, args)
	if err != nil {
		fmt.Println(err)
	}
	// if retries were configured, use error to enqueue a retry, etc.
}
