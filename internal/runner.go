package worker

import (
	"context"
	"fmt"
	"time"
)

type Runner interface {
	DoTask(task Task)
	MaxDuration() time.Duration
}

type RunnerImpl struct {
	manager Manager
}

// could run this as `go r.DoTask(task)`
func (ri *RunnerImpl) DoTask(task Task) {
	timeout := task.Timeout()

	ctx := ri.manager.Context()

	if timeout != nil {
		ctx, cancelFn := context.WithTimeout(ctx, *timeout)
		ctx = ctx
		defer cancelFn()
	}

	args := task.Args() // assume that args are the right type

	name := task.Name()
	taskFn := ri.manager.GetRegisteredTaskFunction(name) // assume taskFn is the right type

	// TODO: do something with this error
	err := taskFn(ctx, args)

	fmt.Println(err)
	// if retries were configured, use error to enqueue a retry, etc.
}

/*
user writes
func(ctx context.Context, an arg of arbitrary type) error
serializer
deserializer
*/

type JobOpt string

func Enqueue(ctx context.Context, taskName string, serializedArgs string, opts ...JobOpt)
