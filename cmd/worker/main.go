package main

import (
	"context"
	"fmt"

	"github.com/sethbacaner/work/gen"
	"github.com/sethbacaner/work/internal"
	"github.com/sethbacaner/work/internal/queue"
)

func main() {

	queue := queue.NewInMemoryQueue()

	manager := internal.NewManager(internal.WithQueue(queue)) // inject test queue into manager
	gen.RegisterAdder(manager, &AdderImpl{})

	enqueueTasks(queue, manager.Registry())

	manager.Start()
}

func enqueueTasks(queue queue.Queue, registry *internal.Registry) {
	ctx := context.Background()
	for i := 0; i < 10; i++ {
		args := &gen.AdderArgs{
			A: 1,
			B: 2,
		}

		serializedTask, _, err := internal.Serialize("Adder", args, *registry)
		if err != nil {
			fmt.Println(err)
			return
		}

		queue.Enqueue(ctx, serializedTask)
	}
}

type AdderImpl struct{}

func (a *AdderImpl) Invoke(ctx context.Context, args gen.AdderArgs) error {
	fmt.Println(args.A + args.B)
	return nil
}
