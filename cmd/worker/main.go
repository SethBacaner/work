package main

import (
	"context"
	"fmt"
	worker "github.com/samuelbacaner/worker/internal"
)

/*
Just finished pair session on 6/20/21

For next time:
1. add functional options to Manager New()
https://www.sohamkamani.com/golang/options-pattern/
2. inject a TestQueue into Manager
3. implement goroutine which writes an AdderTask to TestQueue once a second for 10 seconds
4. run it and see what's broken


Sam:
did a little bit of research on templating and implemented basic code generation with some copy pasta: https://www.calhoun.io/using-code-generation-to-survive-without-generics-in-go/

Generating the args struct is a little tricker.
How about instead of having the user provide a yaml or json definition, we just have them point to a file which has an actual
struct that already has the json tags? That way we'll just suck that up and use it directly. We could do some kind of basic
testing to make sure it's serializable, but basically this means it's on the user to provide a serializable struct.
We could extend it to custom serializers eventually, but let's see how far this takes us.

As to how to parse a golang file to find a specific struct/interface, gomock is one library that does this that we
can look to for inspiration


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
