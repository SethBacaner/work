package queue

import (
	"context"
	"errors"
)

var EmptyQueue = errors.New("no values present in queue")

type Queue interface {
	Enqueue(ctx context.Context, value string) error
	Dequeue(ctx context.Context) (string, error)
}
