package internal

import (
	"context"
)

type TaskFunction func(ctx context.Context, args interface{}) error

type ArgsFactory func() interface{}

type TaskDefinition struct {
	Name        string
	TaskFn      TaskFunction
	argsFactory ArgsFactory
}
