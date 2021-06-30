package queue

import (
	"context"
	"sync"
)

type InMemoryQueue struct {
	values []string
	mu     *sync.Mutex
}

func NewInMemoryQueue() Queue {
	return &InMemoryQueue{
		values: make([]string, 0),
		mu:     &sync.Mutex{},
	}
}

func (im *InMemoryQueue) Enqueue(ctx context.Context, value string) error {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.values = append(im.values, value)

	return nil
}

func (im *InMemoryQueue) Dequeue(ctx context.Context) (string, error) {
	im.mu.Lock()
	defer im.mu.Unlock()

	if len(im.values) == 0 {
		return "", EmptyQueue
	}

	value := im.values[0]
	im.values = im.values[1:len(im.values)]

	return value, nil
}
