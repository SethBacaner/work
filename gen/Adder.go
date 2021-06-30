package gen

import (
	"context"
	"github.com/sethbacaner/work/internal"
)

type AdderArgs struct {
	A int `json:"a"`
	B int `json:"b"`
}

type Adder interface {
	Invoke(ctx context.Context, args AdderArgs) error
}

func RegisterAdder(manager internal.Manager, Adder Adder) {

	taskFn := func(ctx context.Context, args interface{}) error {
		AdderArgs, ok := args.(AdderArgs)
		if !ok {
			// TODO: we need to blow up aggresively
		}

		return Adder.Invoke(ctx, AdderArgs)
	}

	argsFactory := func() interface{} {
		return &AdderArgs{}
	}

	manager.RegisterTask("Adder", argsFactory, taskFn)
}

// TODO: need to generate EnqueuerCLient
// TODO: need to generate serializer and deserializer
