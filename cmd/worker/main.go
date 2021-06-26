package main

import (
	"context"
	"fmt"

	worker "github.com/sethbacaner/worker/internal"
)

/*
Just finished pair session on 6/20/21

For next time:
1. add functional options to Manager New()
https://www.sohamkamani.com/golang/options-pattern/
2. inject a TestQueue into Manager
3. implement goroutine which writes an AdderTask to TestQueue once a second for 10 seconds
4. run it and see what's broken
*/

func main() {

	manager := worker.New() // inject test queue into manager
	RegisterAdder(manager, &AdderImpl{})
	manager.Start()
}

// generated from task definition. tbd on how to define (yaml, etc.)

type Adder interface {
	Invoke(ctx context.Context, args AdderArgs) error
}

type AdderArgs struct {
	a int
	b int
}

func RegisterAdder(manager worker.Manager, adder Adder) {

	taskFn := func(ctx context.Context, args interface{}) error {
		adderArgs, ok := args.(AdderArgs)
		if !ok {
			// TODO: we need to blow up aggresively
		}

		return adder.Invoke(ctx, adderArgs)
	}

	manager.RegisterTask("adder", taskFn)
}

// implemented

type AdderImpl struct{}

func (a *AdderImpl) Invoke(ctx context.Context, args AdderArgs) error {
	fmt.Println(args.a + args.b)
	return nil
}
