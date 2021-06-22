package worker

import (
	"context"
	"errors"
	"sync"
)

var NoTasksEnqueued = errors.New("no tasks present in queue")

type Queue interface {
	EnqueueTask(ctx context.Context, task *Task) error
	GetTask(ctx context.Context) (*Task, error)
}

// TODO: manufactor some kind of task every X seconds
type TestQueue struct {
	tasks []*Task
	mu    *sync.Mutex
}

func (tq *TestQueue) EnqueueTask(ctx context.Context, task *Task) error {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	tq.tasks = append(tq.tasks, task)

	return nil
}

func (tq *TestQueue) GetTask(ctx context.Context) (*Task, error) {
	tq.mu.Lock()
	defer tq.mu.Unlock()

	if len(tq.tasks) == 0 {
		return nil, NoTasksEnqueued
	}

	task := tq.tasks[0]
	tq.tasks = tq.tasks[1:len(tq.tasks)]

	return task, nil
}
