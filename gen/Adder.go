package gen

import (
	"context"
	"github.com/samuelbacaner/worker/internal"
)

// TODO: this needs to be generated more manually
type AdderArgs struct {
	a int
	b int
}

type Adder interface {
	Invoke(ctx context.Context, args AdderArgs) error
}

func RegisterAdder(manager worker.Manager, Adder Adder) {

	taskFn := func(ctx context.Context, args interface{}) error {
		AdderArgs, ok := args.(AdderArgs)
		if !ok {
			// TODO: we need to blow up aggresively
		}

		return Adder.Invoke(ctx, AdderArgs)
	}

	manager.RegisterTask("Adder", taskFn)
}

// TODO: need to generate EnqueuerCLient
// TODO: need to generate serializer and deserializer
